package utilties

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	switch os.Args[1] {
	case "hmac":
		generateAuthHeader()

	}
}

func generateAuthHeader() {
	endpoint := os.Args[2]
	secret := os.Args[3]

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
