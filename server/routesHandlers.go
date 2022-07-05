/**
 *  Authored: Josh Spicer <hello@joshspicer.com>
 */

package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.String(http.StatusOK, "healthy")
}

func TrustedHmacAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		nonce := c.Request.Header.Get("X-Jarvis-Timestamp")

		if authHeader == "" || nonce == "" {
			c.AbortWithStatus(401)
		}

		var timestamp time.Time

		splitted := strings.Split(nonce, "_")
		if len(splitted) != 2 {
			c.AbortWithStatus(401)
		}

		timestampAsInt, err := strconv.ParseInt(splitted[0], 10, 64)
		if err != nil {
			fmt.Printf("Failed to parse timestamp from body: %s\n", strings.ReplaceAll(splitted[0], "\n", ""))
			c.AbortWithStatus(401)
		}

		timestamp = time.Unix(timestampAsInt, 0)

		if c.IsAborted() {
			return
		}

		// The nonce encodes 'timestamp'.
		// Do not accept requests
		// with a 'timestamp' a +/- 1 minute from system time.
		currentTime := time.Now()
		timeDiff := currentTime.Sub(timestamp).Seconds()
		fmt.Printf("Time difference %f seconds\n", timeDiff)
		if math.Abs(timeDiff) > 60 {
			fmt.Printf("Absolute time difference (%f seconds), greater than threshold!\n", timeDiff)
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

		actorLen := len(trustedActors)
		fmt.Printf("Checking received hash against '%d' trusted actors\n", actorLen)

		// Regenerate hash for each trusted actor and compare.
		for i := 0; i < actorLen; i++ {
			actor := trustedActors[i]
			fmt.Printf("Checking hash against actor %s....\n", actor.name)

			h := hmac.New(
				sha256.New,
				[]byte(actor.secret))
			h.Write([]byte(nonce))
			computedHash := hex.EncodeToString(h.Sum(nil))

			// fmt.Printf("\nsecret: %s \ncomputed: %s \nexpected: %s \n\n", actor.secret, computedHash, authHeader)

			if computedHash == authHeader {
				matchStr := fmt.Sprintf("✅ Hash match: %s\n", actor.name)
				fmt.Println(matchStr)
				// bot.SendMessageToPrimaryTelegramGroup(matchStr)
				c.Set("authenticatedUser", actor.name)

				device := c.Request.Header.Get("X-Jarvis-Device")
				fmt.Printf("Device: %s\n", strings.ReplaceAll(device, "\n", ""))

				return
			}
		}
		// Fallback
		bot.SendMessageToPrimaryTelegramGroup("⚠️ Invalid authentication hash provided.")
		c.AbortWithStatus(401)
	}
}

func TrustedKnock(c *gin.Context) {
	// Protected by 'TrustedHmacAuthentication' middleware
	authenticatedUser := c.MustGet("authenticatedUser").(string)

	bot := c.MustGet(BOT_CONTEXT).(*BotExtended)

	matchStr := fmt.Sprintf("✅ %s has knocked.\n", authenticatedUser)
	bot.SendMessageToPrimaryTelegramGroup(matchStr)

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

func Welcome(c *gin.Context) {

	bot := c.MustGet(BOT_CONTEXT).(*BotExtended)
	invite_code := c.Param("invite_code")

	// TODO
	bot.SendMessageToPrimaryTelegramGroup(fmt.Sprintf("Welcome %s", invite_code))
	c.String(http.StatusAccepted, "Welcome, "+invite_code)
}

func Telemetry(c *gin.Context) {
	// Protected by 'TrustedHmacAuthentication' middleware
	authenticatedUser := c.MustGet("authenticatedUser").(string)
	fmt.Printf("Parsing posted telemetry for %s\n", authenticatedUser)

	// Parse request Body
	var telemetry TelemetryPayload
	var bodyBytes []byte
	var err error
	if c.Request.Body != nil {
		bodyBytes, err = ioutil.ReadAll(c.Request.Body)
		if len(bodyBytes) == 0 {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		// Parse JSON into TelemetryPayload
		err = json.Unmarshal(bodyBytes, &telemetry)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	if c.IsAborted() {
		return
	}

	fmt.Printf("Status received: %d", telemetry.State)

	// Write telemetry to db
	// TODO: How to get this to work without panicing?
	// az := c.MustGet(AZURE_CONTEXT).(*AzureExtended)

	// NOTE: We partition on the "type" attribute
	containerClient := ConnectToCosmosContainer(false)
	_, err = containerClient.UpsertItem(context.TODO(), azcosmos.NewPartitionKeyString("telemetry"), bodyBytes, nil)
	if err != nil {
		fmt.Printf("Failed : %s\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
