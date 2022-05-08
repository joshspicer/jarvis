package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	router := SetupRouter()
	bot := SetupTelegram()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	log.Printf("Bot authorized on account %s", bot.Self.UserName)

	SetupCommandHandler(bot)

	print(fmt.Sprintf("Serving at http://localhost:%s", PORT))
	router.Run(fmt.Sprintf(":%s", PORT))

}
