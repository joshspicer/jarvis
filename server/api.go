package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const BOT_CONTEXT = "BOT_CONTEXT"
const AUGUST_HTTP_CONTEXT = `AUGUST_HTTP_CONTEXT`

func BotContext(bot *tgbotapi.BotAPI) gin.HandlerFunc {
	return func(c *gin.Context) {

		botExtended := &BotExtended{bot}

		c.Set(BOT_CONTEXT, botExtended)
		c.Next()
	}
}

func AugustHttpClientContext() gin.HandlerFunc {
	return func(c *gin.Context) {

		httpClient := &http.Client{}
		augustHttpClient := &AugustHttpClient{httpClient}

		c.Set(AUGUST_HTTP_CONTEXT, augustHttpClient)
		c.Next()
	}
}

func SetupRouter(bot *tgbotapi.BotAPI) *gin.Engine {
	router := gin.Default()

	router.Use(BotContext(bot))
	router.Use(AugustHttpClientContext())

	// Hello
	router.GET("/", Hello)
	// Health Endpoint
	router.GET("/health", Health)
	// Knocks
	router.POST("/welcome/:invite_code", Welcome)
	router.POST("/trustedknock", TrustedHmacAuthentication(), TrustedKnock)

	return router
}
