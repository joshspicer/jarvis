package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	version = "dev"
	commit  = "0000"
)

func main() {
	fmt.Printf("[jarvis] %s, %s\n", version, commit)

	// Get first argument to determine mode
	mode := determineMode()

	err := godotenv.Load()
	if err != nil {
		log.Printf("NOTE: No .env file loaded\n")
	}

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
	// case "home":
	// 	log.Printf("Starting home")
	// 	bot := SetupTelegram()
	// 	initializeHome(bot, mode)
	// case "node":
	// 	log.Printf("Starting node")
	// 	initializeNode(mode)
	default:
		log.Fatalf("Invalid mode: %s", mode)
	}
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

	sentry.DoHeartbeat()

	// TODO: sentry.Cron()
}

func determineMode() string {
	if len(os.Args) >= 2 {
		return os.Args[1]
	}
	mode := os.Getenv("JARVIS_MODE")
	if mode != "" {
		return mode
	}
	return mode
}
