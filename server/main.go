package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Get first argument to determine mode
	mode := determineMode()

	bot := SetupTelegram()
	log.Printf("Bot authorized on account %s", bot.Self.UserName)

	switch mode {
	case "jarvis":
		log.Printf("Starting Jarvis")
		initializeJarvis(bot, mode)
	case "sentry":
		log.Printf("Starting Sentry")
		initializeSentry(bot, mode)
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

func initializeSentry(bot *tgbotapi.BotAPI, mode string) {

	botExtended := &BotExtended{bot}
	botExtended.SendMessageToPrimaryTelegramGroup("[sentry initializing]")
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
