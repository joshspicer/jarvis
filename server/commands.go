package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func HelpCommand() string {
	return ("- /status\n")
}

func StatusCommand() string {
	return time.Now().Weekday().String()
}

func InviteCommand(args string) string {
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

	expiration := time.Now().Add(time.Hour * 32)

	return fmt.Sprintf("%s has been invited %d times until %d", name, count, expiration)
}
