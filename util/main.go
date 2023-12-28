package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func usageAndExit() {
	fmt.Printf("Usage: %s <mode> <secret> <endpoint>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usageAndExit()
	}

	switch os.Args[1] {
	case "hmac":
		generateAuthHeader()
	default:
		usageAndExit()
	}
}

func generateAuthHeader() {
	if len(os.Args) < 4 {
		usageAndExit()
	}

	secret := os.Args[2]
	endpoint := os.Args[3]

	timestamp := time.Now().Unix()
	uuid, _ := uuid.NewRandom()
	nonce := fmt.Sprintf("%d_%s", timestamp, uuid)

	h := hmac.New(
		sha256.New,
		[]byte(secret))
	h.Write([]byte(nonce))

	auth := hex.EncodeToString(h.Sum(nil))

	fmt.Printf("%s\n", nonce)
	fmt.Printf("%s\n", auth)
	fmt.Println()
	fmt.Printf("curl http://127.0.0.1:4000/%s -H 'Authorization: %s' -H 'X-Jarvis-Timestamp: %s'\n", endpoint, auth, nonce)
}
