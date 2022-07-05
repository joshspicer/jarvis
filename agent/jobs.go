package main

import (
	"fmt"
	"net/http"
)

func Telemetry(httpClient http.Client, serviceUri string) {
	// Create HTTP Client

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/telemetry", serviceUri), nil)
}
