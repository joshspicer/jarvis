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

		// Retrieve list of trusted actors
		var trustedActors []Actor = GetTrustedActors()
		taLen := len(trustedActors)
		if taLen == 0 {
			println("No trusted actors available.")
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		// Regenerate hash for each trusted actor and compare.
		for i := 0; i < taLen; i++ {
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

		// Fallback
		c.AbortWithStatus(401)
	}
}

func TrustedKnock(c *gin.Context) {
	c.String(http.StatusAccepted, "trusted knock")
}
