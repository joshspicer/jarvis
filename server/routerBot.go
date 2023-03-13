package main

func RouterCommandHandler(bot *BotExtended, command string, args string) string {
	// Extract the command from the Message.
	switch command {
	case "help":
		return "router help!"
	case "status":
		return "router status!"
	default:
		return "Try Again."
	}
}
