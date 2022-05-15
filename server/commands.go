package main

import (
	"time"
)

func HelpCommand() string {
	return ("- /status\n")
}

func StatusCommand() string {
	return time.Now().Weekday().String()
}
