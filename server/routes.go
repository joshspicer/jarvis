package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	c.String(http.StatusOK, "hello")
}

func Health(c *gin.Context) {
	c.String(http.StatusOK, "healthy")
}

func TrustedHmacAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(401)
		}

		var nonce string
		var timestamp time.Time

		// Parse request Body
		if c.Request.Body != nil {

			bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
			if len(bodyBytes) == 0 {
				c.AbortWithStatus(401)
			}

			nonce = string(bodyBytes)

			splitted := strings.Split(nonce, "_")
			if len(splitted) != 2 {
				c.AbortWithStatus(401)
			}

			timestampAsInt, err := strconv.ParseInt(splitted[0], 10, 64)
			if err != nil {
				fmt.Printf("Failed to parse timestamp from body: %s\n", splitted[0])
				c.AbortWithStatus(401)
			}

			timestamp = time.Unix(timestampAsInt, 0)
		}

		if c.IsAborted() {
			return
		}

		// The nonce encodes 'timestamp'.
		// Do not accept requests
		// with a 'timestamp' a +/- 1 minute from system time.
		currentTime := time.Now()
		timeDiff := currentTime.Sub(timestamp).Seconds()
		fmt.Printf("Time difference %f seconds\n", timeDiff)
		if timeDiff > 60 {
			fmt.Printf("Time difference %f seconds\n", timeDiff)
			c.AbortWithStatus(401)
		}

		if c.IsAborted() {
			return
		}

		bot := c.MustGet(BOT_CONTEXT).(*BotExtended)

		// Retrieve list of trusted actors

		trustedActors, err := GetTrustedActors()
		if err != nil {
			fmt.Printf("Failed to retrieve trusted actors: %s\n", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		// Regenerate hash for each trusted actor and compare.
		for i := 0; i < len(trustedActors); i++ {
			actor := trustedActors[i]

			h := hmac.New(
				sha256.New,
				[]byte(actor.secret))
			h.Write([]byte(nonce))
			computedHash := hex.EncodeToString(h.Sum(nil))

			// fmt.Printf("\nsecret: %s \ncomputed: %s \nexpected: %s \n\n", actor.secret, computedHash, authHeader)

			if computedHash == authHeader {
				fmt.Printf("Hash match: %s\n", actor.name)
				c.String(http.StatusAccepted, actor.name)
				return
			}
		}

		bot.SendMessageToPrimaryTelegramGroup("[!] An attempt to validate an HMAC hash failed.")

		// Fallback
		c.AbortWithStatus(401)
	}
}

func TrustedKnock(c *gin.Context) {

	august := c.MustGet(AUGUST_HTTP_CONTEXT).(*AugustHttpClient)

	error := august.OperateLock("unlock")
	if error != nil {
		fmt.Println(fmt.Errorf("failed to unlock August: %s", error))
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	// Accept if we have not aborted.
	if !c.IsAborted() {
		c.String(http.StatusAccepted, "trusted knock")
	}
}

func Welcome(c *gin.Context) {

	bot := c.MustGet(BOT_CONTEXT).(*BotExtended)
	invite_code := c.Param("invite_code")

	bot.SendMessageToPrimaryTelegramGroup(fmt.Sprintf("Welcome %s", invite_code))

	c.String(http.StatusAccepted, "welcome, "+invite_code)
}
