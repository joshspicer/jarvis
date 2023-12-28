package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateAuthHeaderForPrimaryActor() (string, string, error) {

	trustedActors, err := GetTrustedActors()
	if err != nil {
		fmt.Printf("Failed to retrieve trusted actors: %s\n", err)
		return "", "", err
	}
	var primaryActor = trustedActors[0]

	timestamp := time.Now().Unix()
	uuid, _ := uuid.NewRandom()
	nonce := fmt.Sprintf("%d_%s", timestamp, uuid)

	h := hmac.New(
		sha256.New,
		[]byte(primaryActor.secret))
	h.Write([]byte(nonce))

	auth := hex.EncodeToString(h.Sum(nil))

	return auth, nonce, nil
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		nonce := c.Request.Header.Get("X-Jarvis-Timestamp")

		if authHeader == "" || nonce == "" {
			c.AbortWithStatus(401)
			return
		}

		var timestamp time.Time

		splitted := strings.Split(nonce, "_")
		if len(splitted) != 2 {
			c.AbortWithStatus(401)
			return
		}

		timestampAsInt, err := strconv.ParseInt(splitted[0], 10, 64)
		if err != nil {
			fmt.Printf("Failed to parse timestamp from body: %s\n", strings.ReplaceAll(splitted[0], "\n", ""))
			c.AbortWithStatus(401)
			return
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

		fetchedBot, hasBot := c.Get(BOT_CONTEXT)
		var bot *BotExtended = nil
		if hasBot {
			bot = fetchedBot.(*BotExtended)
		}

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

			computedHash := []byte(hex.EncodeToString(h.Sum(nil)))
			authHeader := []byte(authHeader)

			// fmt.Printf("\nsecret: %s \ncomputed: %s \nexpected: %s \n\n", actor.secret, computedHash, authHeader)

			valid := subtle.ConstantTimeCompare(computedHash, authHeader) == 1 // ConstantTimeCompare returns 1 if the two slices, x and y, have equal contents and 0 otherwise.
			if valid {
				matchStr := fmt.Sprintf("✅ Hash match: %s\n", actor.name)
				fmt.Println(matchStr)
				if hasBot {
					bot.SendMessageToPrimaryTelegramGroup(matchStr)
				}
				c.Set("authenticatedUser", actor.name)
				c.Next()

				device := c.Request.Header.Get("X-Jarvis-Device")
				fmt.Printf("Device: %s\n", strings.ReplaceAll(device, "\n", ""))

				return
			}
		}

		// Fallback
		if hasBot {
			bot.SendMessageToPrimaryTelegramGroup("⚠️ Invalid authentication hash provided.")
		}
		c.AbortWithStatus(401)
	}
}
