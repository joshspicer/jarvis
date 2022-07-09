/**
 *  Authored: Josh Spicer <hello@joshspicer.com>
 */

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Startup variables set to panic if any errors occur.
	bot := SetupTelegram()

	az := AzureExtended{}
	az.RefreshCredential(true)

	router := SetupRouter(bot, &az)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	log.Printf("Bot authorized on account %s", bot.Self.UserName)

	go SetupCommandHandler(bot)

	print(fmt.Sprintf("Serving at http://localhost:%s", PORT))
	router.Run(fmt.Sprintf(":%s", PORT))

}
