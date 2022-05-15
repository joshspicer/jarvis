package main

import (
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const BOT_CONTEXT = "BOT_CONTEXT"

func BotContext(bot *tgbotapi.BotAPI) gin.HandlerFunc {
	return func(c *gin.Context) {

		botExtended := &BotExtended{bot}

		c.Set(BOT_CONTEXT, botExtended)
		c.Next()
	}
}

func SetupRouter(bot *tgbotapi.BotAPI) *gin.Engine {
	router := gin.Default()

	router.Use(BotContext(bot))

	// Hello
	router.GET("/", Hello)
	// Health Endpoint
	router.GET("/health", Health)
	// Knocks
	router.POST("/welcome/:invite_code", Welcome)
	router.POST("/trustedknock", TrustedHmacAuthentication(), TrustedKnock)

	return router
}
