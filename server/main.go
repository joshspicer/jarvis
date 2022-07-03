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
	azureCredential := InitAzure(true) // Ensures a default credential can be created and checks for existance of required DBs
	router := SetupRouter(bot, azureCredential)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	log.Printf("Bot authorized on account %s", bot.Self.UserName)

	go SetupCommandHandler(bot)

	print(fmt.Sprintf("Serving at http://localhost:%s", PORT))
	router.Run(fmt.Sprintf(":%s", PORT))

}
