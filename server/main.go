package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	version = "dev"
	commit  = "0000"
)

func main() {
	fmt.Printf("[jarvis] %s, %s\n", version, commit)

	err := godotenv.Load()
	if err != nil {
		log.Printf("NOTE: No .env file loaded\n")
	}

	// Get first argument to determine mode
	mode := determineMode()

	_, isRelease := os.LookupEnv("RELEASE")
	if isRelease {
		log.Printf("Mode: Release\n")
		gin.SetMode(gin.ReleaseMode)
	}
	switch mode {
	case "cloud":
		log.Printf("Starting cloud")
		bot := SetupTelegram()
		initializeCloud(bot, mode)
	case "sentry":
		initializeSentry()
	default:
		log.Fatalf("Invalid mode: %s", mode)
	}

	log.Fatalf("At end of main")
}

func initializeCloud(bot *tgbotapi.BotAPI, mode string) {
	router := CloudRouter(bot)
	botExtended := &BotExtended{bot}
	go SetupTelegramCommandHandler(botExtended, mode)

	info := fmt.Sprintf("[jarvis] cloud initializing: %s, %s", version, commit)
	botExtended.SendMessageToPrimaryTelegramGroup(info)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	collection, mongoCtx, cleanup := getMongoCollectionConnection("primary")
	mongoResult := collection.FindOneAndUpdate(mongoCtx, bson.M{"type": "globalState"}, bson.M{"$set": bson.M{"cloudInitializedAt": time.Now().UTC().String()}})

	if mongoResult.Err() != nil {
		log.Fatalf("Failed to update db %s\n", mongoResult.Err())
	}

	cleanup()

	log.Printf("Serving at http://localhost:%s", PORT)
	router.Run(fmt.Sprintf(":%s", PORT))
}

func initializeSentry() {
	jarvisCloudBaseAddr := os.Getenv("JARVIS_CLOUD_BASE_ADDR")
	trustedActors, err := GetTrustedActors()
	if err != nil {
		log.Fatalf("Failed to retrieve trusted actors: %s\n", err)
	}

	primaryActor := trustedActors[0]

	// Check jarvisCloudBaseAddr contains protocol and no trailing slash (eg: https://example.com)
	if !strings.HasPrefix(jarvisCloudBaseAddr, "http") || jarvisCloudBaseAddr[len(jarvisCloudBaseAddr)-1:] == "/" {
		log.Fatalf("Misformed cloud base addr. Got: %s", jarvisCloudBaseAddr)
	}

	sentry := Sentry{
		CloudBaseAddr: jarvisCloudBaseAddr,
		Actor:         primaryActor,
	}

	sentry.DoHeartbeat(INITIALIZE)

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to initialize scheduler: %s\n", err)
	}

	seconds, ok := os.LookupEnv("JARVIS_HEARTBEAT_INTERVAL_SECONDS")
	if !ok {
		seconds = "600"
	}

	// Every X seconds
	_, err = s.NewJob(
		gocron.CronJob(
			fmt.Sprintf("*/%s * * * * *", seconds),
			true,
		),
		gocron.NewTask(
			func() {
				sentry.DoHeartbeat(OK)
			},
		),
	)

	if err != nil {
		log.Fatalf("Failed to initialize job: %s\n", err)
	}

	fmt.Printf("Scheduled heartbeat every '%s' seconds\n", seconds)

	s.Start()

	// Catch an exit signal and run cleanup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		sentry.DoHeartbeat(INTERRUPTED)
		log.Fatalln("\n[jarvis] Interrupted...")
	}()

	// Keep alive
	select {}
}

func determineMode() string {
	if len(os.Args) >= 2 {
		return os.Args[1]
	}
	mode := os.Getenv("JARVIS_MODE")
	fmt.Printf("Mode: %s\n", mode)
	if mode != "" {
		return mode
	}
	return mode
}
