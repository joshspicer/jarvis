package main

import "time"

type StateCode int64

const (
	Healthy StateCode = iota
	GeneralError
)

type TelemetryPayload struct {
	State     StateCode `json:"state"`
	Timestamp time.Time `json:"timestamp"`
}
