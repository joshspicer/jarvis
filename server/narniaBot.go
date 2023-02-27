package main

func NarniaCommandHandler(bot *BotExtended, command string, args string) string {
	// Extract the command from the Message.
	switch command {
	case "help":
		return "narnia Help"
	case "status":
		return "narnia status"
	default:
		return "Try Again."
	}
}
