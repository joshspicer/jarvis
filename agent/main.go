/**
 *  Authored: Josh Spicer <hello@joshspicer.com>
 */

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/robfig/cron/v3"
)

func main() {
	router := SetupAgentRouter()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4001"
	}

	SERVICE_URI := os.Getenv("SERVICE_URI")
	if SERVICE_URI == "" {
		SERVICE_URI = "http://localhost:4000"
	}

	CLIENT_SECRET := os.Getenv("CLIENT_SECRET")
	if CLIENT_SECRET == "" {
		panic("CLIENT_SECRET is not set")
	}

	HEARTBEAT_RATE := os.Getenv("HEARTBEAT_RATE")
	if HEARTBEAT_RATE == "" {
		HEARTBEAT_RATE = "@every 5s"
	}

	httpClient := http.Client{}

	c := cron.New()
	c.AddFunc(HEARTBEAT_RATE, func() {
		Telemetry(httpClient, SERVICE_URI, CLIENT_SECRET)
	})
	c.Start()

	print(fmt.Sprintf("Serving at http://localhost:%s", PORT))
	router.Run(fmt.Sprintf(":%s", PORT))

}
