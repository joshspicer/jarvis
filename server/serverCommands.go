package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ServerHelpCommand() string {
	return ("- /status\n")
}

func ServerStatusCommand() string {
	return time.Now().Weekday().String()
}

func AugustInviteCommand(args string) string {
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
