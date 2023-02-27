package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/autotls"
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

	switch mode {
	case "jarvis":
		log.Printf("Starting Jarvis")
		bot := SetupTelegram()
		initializeJarvis(bot, mode)
	case "narnia":
		log.Printf("Starting narnia")
		bot := SetupTelegram()
		initializeNarnia(bot, mode)
	case "niro":
		log.Printf("Starting niro")
		initializeNiro(mode)
	default:
		log.Fatalf("Invalid mode: %s", mode)
	}
}

func initializeJarvis(bot *tgbotapi.BotAPI, mode string) {
	router := JarvisRouter(bot)
	go SetupTelegramCommandHandler(&BotExtended{bot}, mode)

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
	SetupTelegramCommandHandler(botExtended, mode)
}

func initializeNiro(mode string) {
	router := NiroRouter()

	HTTPS_DOMAINS := os.Getenv("HTTPS_DOMAINS")

	if HTTPS_DOMAINS != "" {
		split := strings.Split(HTTPS_DOMAINS, ",")
		log.Printf("Serving as HTTPS on one of: '%s'", HTTPS_DOMAINS)
		log.Fatal(autotls.Run(router, split...))
	} else {
		PORT := os.Getenv("PORT")
		if PORT == "" {
			PORT = "4000"
		}
		log.Printf("Serving at http://localhost:%s", PORT)
		router.Run(fmt.Sprintf(":%s", PORT))
	}
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
