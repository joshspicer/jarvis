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

	httpClient := http.Client{}

	c := cron.New()
	c.AddFunc("@every 2s", func() {
		Telemetry(httpClient, SERVICE_URI)
	})
	c.Start()

	print(fmt.Sprintf("Serving at http://localhost:%s", PORT))
	router.Run(fmt.Sprintf(":%s", PORT))

}
