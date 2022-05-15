package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	bot := SetupTelegram()
	router := SetupRouter(bot)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	log.Printf("Bot authorized on account %s", bot.Self.UserName)

	go SetupCommandHandler(bot)

	print(fmt.Sprintf("Serving at http://localhost:%s", PORT))
	router.Run(fmt.Sprintf(":%s", PORT))

}
