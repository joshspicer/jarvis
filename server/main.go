package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Get first argument to determine mode
	mode := determineMode()

	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file loaded")
	}

	bot := SetupTelegram()
	log.Printf("Bot authorized on account %s", bot.Self.UserName)

	switch mode {
	case "jarvis":
		log.Printf("Starting Jarvis")
		initializeJarvis(bot, mode)
	case "narnia":
		log.Printf("Starting narnia")
		initializeNarnia(bot, mode)
	default:
		log.Fatalf("Invalid mode: %s", mode)
	}
}

func initializeJarvis(bot *tgbotapi.BotAPI, mode string) {

	router := SetupServerRouter(bot)
	go SetupCommandHandler(&BotExtended{bot}, mode)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	log.Printf("Serving at http://localhost:%s", PORT)
	router.Run(fmt.Sprintf(":%s", PORT))
}

func initializeNarnia(bot *tgbotapi.BotAPI, mode string) {

	botExtended := &BotExtended{bot}
	botExtended.SendMessageToPrimaryTelegramGroup("[narnia initializing]")
	SetupCommandHandler(botExtended, mode)
}

func determineMode() string {
	DEFAULT_MODE := "jarvis"
	if len(os.Args) >= 2 {
		return os.Args[1]
	}
	mode := os.Getenv("JARVIS_MODE")
	if mode != "" {
		return mode
	}
	return DEFAULT_MODE
}
