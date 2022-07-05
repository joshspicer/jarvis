// Client

package main

import "time"

type StateCode int64

const (
	Healthy StateCode = iota
	GeneralError
)

type TelemetryPayload struct {
	Id        string    `json:"id"`
	Type      string    `json:"type"`
	State     StateCode `json:"state"`
	Timestamp time.Time `json:"timestamp"`
}
