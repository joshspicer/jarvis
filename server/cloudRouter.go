package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CloudRouter(bot *tgbotapi.BotAPI) *gin.Engine {
	router := gin.Default()

	router.Use(BotContext(bot))
	router.Use(AugustHttpClientContext())

	// Static
	router.StaticFile("robots.txt", "./static/robots.txt")

	// Meta
	router.GET("/health", health)
	router.GET("/whoami", Auth(), whoami)

	router.POST("/heartbeat", Auth(), sentryHeartbeatHandler)

	// Knocks
	// router.POST("/welcome/:invite_code", welcome)
	// router.POST("/trustedknock", Auth(), trustedKnock)

	return router
}

func health(c *gin.Context) {
	// Set no cache headers
	c.Header("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("X-Accel-Expires", "0")

	c.String(http.StatusOK, fmt.Sprintf("healthy\n%s\n%s", version, commit))
}

func whoami(c *gin.Context) {
	// Protected by 'TrustedHmacAuthentication' middleware
	authenticatedUser := c.MustGet("authenticatedUser").(string)

	// Set no cache headers
	c.Header("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("X-Accel-Expires", "0")

	c.String(http.StatusOK, authenticatedUser)
}

func sentryHeartbeatHandler(c *gin.Context) {
	// Protected by 'TrustedHmacAuthentication' middleware
	authenticatedUser := c.MustGet("authenticatedUser").(string)

	// Set no cache headers
	c.Header("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Header("X-Accel-Expires", "0")

	// Parse JSON
	var heartbeat Heartbeat
	if err := c.ShouldBindJSON(&heartbeat); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if heartbeat.Code > 0 {
		bot := c.MustGet(BOT_CONTEXT).(*BotExtended)
		bot.SendMessageToPrimaryTelegramGroup(fmt.Sprintf("Heartbeat from '%s' on '%s' at %d", authenticatedUser, heartbeat.HostName, heartbeat.Timestamp))
	}

	// Accept if we have not aborted.
	if !c.IsAborted() {
		var response HeartbeatResponse
		response.Accepted = true
		c.JSON(http.StatusAccepted, response)
		return
	}

	// Reject if down here
	c.Status(http.StatusInternalServerError)
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
