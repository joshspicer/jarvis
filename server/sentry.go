package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Sentry struct {
	CloudBaseAddr string
	Actor         Actor
}

// Heartbeat object definition (JSON)
type Heartbeat struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	HostName  string `json:"hostname"`
}

type HeartbeatResponse struct {
	accepted bool
}

func (s Sentry) DoHeartbeat() {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "(unknown)"
	}

	values := Heartbeat{
		Id:        s.Actor.name,
		Timestamp: time.Now().Unix(),
		HostName:  hostName,
	}
	s.sendHeartbeat(values)
}

func (s Sentry) sendHeartbeat(values Heartbeat) {

	json_data, err := json.Marshal(values)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Create new POST request object
	req, err := http.NewRequest("POST", s.CloudBaseAddr+"/heartbeat", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println(err.Error())
	}

	auth, nonce, err := GenerateAuthHeaderForPrimaryActor()
	if err != nil {
		fmt.Println(err.Error())
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)
	req.Header.Set("X-Jarvis-Timestamp", nonce)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Close request
	defer resp.Body.Close()

	// Parse response
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	log.Println(res["json"])
}
