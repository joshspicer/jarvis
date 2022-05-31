package main

import (
	"encoding/json"
	"net/http"
	"os"
)

const AUGUST_SESSION_URL = "https://api-production.august.com/session"
const AUGUST_LOCK_UNLOCK_URL = "blah"

type AugustHttpClient struct {
	*http.Client
}

type AugustSessionPayload struct {
	ApiKey    string `json:"apiKey"`
	InstallId string `json:"installID"`
	Password  string `json:"password"`
	IdType    string `json:"IDType"`
	AugustId  string `json:"augustID"`
}

func getAugustSessionEnvVars() (string, error) {

	augustSessionPayload := AugustSessionPayload{
		ApiKey:    os.Getenv("AUGUST_API_KEY"),
		InstallId: os.Getenv("AUGUST_INSTALLID"),
		Password:  os.Getenv("AUGUST_PASSWORD"),
		IdType:    os.Getenv("AUGUST_ID_TYPE"),
		AugustId:  os.Getenv("AUGUST_ID"),
	}

	json, err := json.Marshal(augustSessionPayload)

	return string(json), nil
}

func (client AugustHttpClient) getAugustSession() (string, error) {
	return "nice", nil
}

func (client AugustHttpClient) Unlock() {

	session, err := client.getAugustSession()

}
