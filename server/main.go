package main

import (
	"fmt"
	"log"
	"os"

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
	case "cluster":
		log.Printf("Starting cluster")
		bot := SetupTelegram()
		initializeCluster(bot, mode)
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

func initializeCluster(bot *tgbotapi.BotAPI, mode string) {
	router := ClusterRouter(bot)
	botExtended := &BotExtended{bot}
	go SetupTelegramCommandHandler(botExtended, mode)

	info := fmt.Sprintf("[jarvis] cluster initializing: %s, %s", version, commit)
	botExtended.SendMessageToPrimaryTelegramGroup(info)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	log.Printf("Serving at http://localhost:%s", PORT)
	router.Run(fmt.Sprintf(":%s", PORT))
}

// func initializeHome(bot *tgbotapi.BotAPI, mode string) {
// 	botExtended := &BotExtended{bot}
// 	botExtended.SendMessageToPrimaryTelegramGroup("[narnia initializing]")
// 	SetupTelegramCommandHandler(botExtended, mode)
// }

// func initializeNode(mode string) {
// 	router := NodeRouter()

// 	PORT := os.Getenv("PORT")
// 	if PORT == "" {
// 		PORT = "4000"
// 	}
// 	log.Printf("Serving at http://localhost:%s", PORT)
// 	router.Run(fmt.Sprintf(":%s", PORT))

// }

func determineMode() string {
	DEFAULT_MODE := "cluster"
	if len(os.Args) >= 2 {
		return os.Args[1]
	}
	mode := os.Getenv("JARVIS_MODE")
	if mode != "" {
		return mode
	}
	return DEFAULT_MODE
}
