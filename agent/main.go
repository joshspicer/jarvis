/**
 *  Authored: Josh Spicer <hello@joshspicer.com>
 */

package main

import (
	"fmt"
	"os"
)

func main() {
	router := SetupAgentRouter()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	print(fmt.Sprintf("Serving at http://localhost:%s", PORT))
	router.Run(fmt.Sprintf(":%s", PORT))

}
