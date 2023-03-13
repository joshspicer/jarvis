package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
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
	case "router":
		log.Printf("Starting router")
		bot := SetupTelegram()
		initializeRouter(bot, mode)
	case "node":
		log.Printf("Starting node")
		initializeNode(mode)
	default:
		log.Fatalf("Invalid mode: %s", mode)
	}
}

func initializeCluster(bot *tgbotapi.BotAPI, mode string) {
	router := ClusterRouter(bot)
	go SetupTelegramCommandHandler(&BotExtended{bot}, mode)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	log.Printf("Serving at http://localhost:%s", PORT)
	router.Run(fmt.Sprintf(":%s", PORT))
}

func initializeRouter(bot *tgbotapi.BotAPI, mode string) {
	botExtended := &BotExtended{bot}
	botExtended.SendMessageToPrimaryTelegramGroup("[narnia initializing]")
	SetupTelegramCommandHandler(botExtended, mode)
}

func initializeNode(mode string) {
	router := NodeRouter()

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
