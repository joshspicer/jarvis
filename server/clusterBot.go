package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ClusterCommandHandler(bot *BotExtended, command string, args string) string {
	// Extract the command from the Message.
	switch command {
	case "help":
		return clustersHelpCommand()
	case "status":
		return clusterStatusCommand()
	case "invite":
		return augustInviteCommand(args)
	default:
		return "Try Again."
	}
}

func clustersHelpCommand() string {
	return ("- /status\n")
}

func clusterStatusCommand() string {
	return time.Now().Weekday().String()
}

func augustInviteCommand(args string) string {
	split := strings.Split(args, " ")
	if len(split) > 2 {
		return "Too many arguments. Usage: /invite <name> [count=1]"
	}
	name := split[0]
	var count int64 = 1
	if len(split) == 2 {
		providedCount, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			return "Invalid count. Must be integer"
		}
		count = providedCount
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return "Failed to generate unique invite Id"
	}
	expiration := time.Now().Add(time.Hour * 32)

	// TODO: Acually do it.

	return fmt.Sprintf("%s has been invited %d times until %s.  Invite code: %s", name, count, expiration.Format(time.RFC822Z), uuid.String())
}
