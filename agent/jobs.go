package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func Telemetry(httpClient http.Client, serviceUri string, clientSecret string) {

	payload, host := calculatePayload()
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/telemetry", serviceUri), payload)

	now := time.Now()
	nonce := fmt.Sprintf("%d_%s", now.Unix(), uuid.New().String())
	h := hmac.New(
		sha256.New,
		[]byte(clientSecret))
	h.Write([]byte(nonce))
	computedHash := hex.EncodeToString(h.Sum(nil))

	req.Header.Set("Authorization", computedHash)
	req.Header.Add("X-Jarvis-Timestamp", nonce)
	req.Header.Add("X-Jarvis-Device", host)

	res, _ := httpClient.Do(req)

	fmt.Printf("[Telemetry] %s - %d\n", now.String(), res.StatusCode)
}

func calculatePayload() (io.Reader, string) {
	payload := TelemetryPayload{}
	payload.State = 0
	payload.Timestamp = time.Now()
	payload.Type = "telemetry"

	// Id is the host running the agent
	hostname, _ := os.Hostname()
	if hostname != "" {
		payload.Id = hostname
	} else {
		payload.Id = os.Getenv("HOSTNAME")
		if payload.Id == "" {
			payload.Id = "Generic Host"
		}
	}

	myBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal payload: %s\n", err.Error())
	}

	return bytes.NewReader(myBytes), payload.Id
}
