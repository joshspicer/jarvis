package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ClusterRouter(bot *tgbotapi.BotAPI) *gin.Engine {
	router := gin.Default()

	router.Use(BotContext(bot))
	router.Use(AugustHttpClientContext())

	// Static
	router.StaticFile("robots.txt", "./static/robots.txt")

	router.GET("/health", health)
	// Knocks
	router.POST("/welcome/:invite_code", welcome)
	router.POST("/trustedknock", Auth(), trustedKnock)

	return router
}

func health(c *gin.Context) {
	c.String(http.StatusOK, "healthy")
}

func trustedKnock(c *gin.Context) {
	// Protected by 'TrustedHmacAuthentication' middleware
	authenticatedUser := c.MustGet("authenticatedUser").(string)

	august := c.MustGet(AUGUST_HTTP_CONTEXT).(*AugustHttpClient)

	error := august.OperateLock("unlock")
	if error != nil {
		fmt.Println(fmt.Errorf("failed to unlock August: %s", error))
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	// Accept if we have not aborted.
	if !c.IsAborted() {
		c.String(http.StatusAccepted, fmt.Sprintf("Welcome, %s.", authenticatedUser))
	}
}

func welcome(c *gin.Context) {

	bot := c.MustGet(BOT_CONTEXT).(*BotExtended)
	invite_code := c.Param("invite_code")

	// TODO
	bot.SendMessageToPrimaryTelegramGroup(fmt.Sprintf("Welcome %s", invite_code))
	c.String(http.StatusAccepted, "Welcome, "+invite_code)
}
