package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const BOT_CONTEXT = "BOT_CONTEXT"

type BotExtended struct {
	*tgbotapi.BotAPI
}

func BotContext(bot *tgbotapi.BotAPI) gin.HandlerFunc {
	return func(c *gin.Context) {

		botExtended := &BotExtended{bot}

		c.Set(BOT_CONTEXT, botExtended)
		c.Next()
	}
}

func (b *BotExtended) SendMessageToPrimaryTelegramGroup(message string) {
	// Get primary group, which is the first in the space-separated list.
	validTelegramGroups := strings.Split(os.Getenv("VALID_TELEGRAM_GROUPS"), " ")

	if len(validTelegramGroups) == 0 {
		log.Panic("No valid Telegram groups configured.")
	}

	primaryChatId, err := strconv.ParseInt(validTelegramGroups[0], 10, 64)
	if err != nil {
		log.Panic(err)
	}

	msg := tgbotapi.NewMessage(primaryChatId, message)
	b.Send(msg)
}

func SetupTelegram() *tgbotapi.BotAPI {
	TELEGRAM_BOT_TOKEN := os.Getenv("TELEGRAM_BOT_TOKEN")
	if TELEGRAM_BOT_TOKEN == "" {
		panic("TELEGRAM_BOT_TOKEN is required.")
	}

	bot, err := tgbotapi.NewBotAPI(TELEGRAM_BOT_TOKEN)
	if err != nil {
		panic(fmt.Sprintf("Error Creating new Telegram bot object: %s", err))
	}

	log.Printf("Bot authorized on account %s", bot.Self.UserName)

	return bot
}

func validate(configEnvVariable string, valueToCheck string) bool {

	validString := os.Getenv(configEnvVariable)
	if validString == "" {
		return false
	}
	valid := strings.Split(validString, " ")
	length := len(valid)

	for i := 0; i < length; i++ {
		if valueToCheck == valid[i] {
			return true
		}
	}

	// Default deny
	return false
}

func SetupTelegramCommandHandler(bot *BotExtended, handlerMode string) {

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		bot.Debug = false
	} else {
		bot.Debug = true
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		isValidSender := validate(
			"VALID_TELEGRAM_SENDERS",
			strconv.FormatInt(update.Message.From.ID, 10))
		if !isValidSender {
			println(fmt.Sprintf("[!] Unauthorized user: %s (%d)", update.Message.From.UserName, update.Message.From.ID))
			continue
		}

		isValidGroup := validate(
			"VALID_TELEGRAM_GROUPS",
			strconv.FormatInt(update.Message.Chat.ID, 10))
		if !isValidGroup {
			title := update.Message.Chat.Title
			sender := update.Message.From.UserName
			println(fmt.Sprintf("[!] Unauthorized chat (type=%s): %s (%s)", update.Message.Chat.Type, sender, title))
			continue
		}

		// ignore any non-Message updates
		if update.Message == nil {
			continue
		}

		// ignore any non-command Messages
		if !update.Message.IsCommand() {
			continue
		}

		// Create a new MessageConfig.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		command := update.Message.Command()
		args := update.Message.CommandArguments()
		switch handlerMode {
		case "jarvis":
			msg.Text = JarvisCommandHandler(bot, command, args)
		case "narnia":
			msg.Text = NarniaCommandHandler(bot, command, args)
		default:
			msg.Text = "[ERR] Invalid handler mode!"
			log.Printf("Invalid handler mode %s", handlerMode)

		}

		// Send result as response to the user
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
		}
	}
}
