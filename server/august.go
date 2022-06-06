package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const AUGUST_SESSION_URL = "https://api-production.august.com/session"

type AugustHttpClient struct {
	*http.Client
}

type AugustSessionPayload struct {
	ApiKey    string `json:"apiKey"`
	InstallId string `json:"installID"`
	Password  string `json:"password"`
	IdType    string `json:"IDType"`
	AugustId  string `json:"augustID"`
	LockId    string `json:"lockID"`
}

func getAugustSessionEnvVars() (AugustSessionPayload, error) {

	augustSessionPayload := AugustSessionPayload{
		ApiKey:    os.Getenv("AUGUST_API_KEY"),
		InstallId: os.Getenv("AUGUST_INSTALLID"),
		Password:  os.Getenv("AUGUST_PASSWORD"),
		IdType:    os.Getenv("AUGUST_ID_TYPE"),
		AugustId:  os.Getenv("AUGUST_ID"),
		LockId:    os.Getenv("AUGUST_LOCK_ID"),
	}
	return augustSessionPayload, nil
}

func (client AugustHttpClient) getAugustSession() (string, AugustSessionPayload, error) {
	sessionEnvVars, err := getAugustSessionEnvVars()
	if err != nil {
		return "", AugustSessionPayload{}, err
	}

	httpBody, err := json.Marshal(map[string]string{
		"installId":  sessionEnvVars.InstallId,
		"password":   sessionEnvVars.Password,
		"identifier": fmt.Sprintf("%s:%s", sessionEnvVars.IdType, sessionEnvVars.AugustId),
	})
	if err != nil {
		return "", AugustSessionPayload{}, err
	}

	httpRequest, err := http.NewRequest("POST", AUGUST_SESSION_URL, bytes.NewReader(httpBody))
	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Accept-Version", "0.0.1")
	httpRequest.Header.Add("x-august-api-key", sessionEnvVars.ApiKey)
	httpRequest.Header.Add("x-kease-api-key", sessionEnvVars.ApiKey)
	httpRequest.Header.Add("User-Agent", "August/Luna-3.2.2")

	if err != nil {
		return "", AugustSessionPayload{}, fmt.Errorf("error creating request: %v", err)
	}

	res, err := client.Do(httpRequest)
	if err != nil || res.StatusCode > 299 {
		return "", AugustSessionPayload{}, fmt.Errorf("error sending request (code=%d): %v", res.StatusCode, err)
	}

	accessToken := res.Header.Get("x-august-access-token")
	if accessToken == "" {
		return "", AugustSessionPayload{}, errors.New("response doesn't contain a non-empty access token")
	}

	return accessToken, sessionEnvVars, nil
}

func (client AugustHttpClient) OperateLock(mode string) error {

	if mode != "lock" && mode != "unlock" {
		return errors.New("mode must be either 'lock' or 'unlock'")
	}

	session, sessionEnvVars, err := client.getAugustSession()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api-production.august.com/remoteoperate/%s/%s", sessionEnvVars.LockId, mode)

	httpRequest, err := http.NewRequest("PUT", url, strings.NewReader(""))

	httpRequest.Header.Add("Content-Type", "application/json")
	httpRequest.Header.Add("Content-Length", "0")
	httpRequest.Header.Add("Accept-Version", "0.0.1")
	httpRequest.Header.Add("x-august-api-key", sessionEnvVars.ApiKey)
	httpRequest.Header.Add("x-kease-api-key", sessionEnvVars.ApiKey)
	httpRequest.Header.Add("User-Agent", "August/Luna-3.2.2")
	httpRequest.Header.Add("x-august-access-token", session)

	if err != nil {
		return err
	}

	res, err := client.Do(httpRequest)

	if err != nil || res.StatusCode > 299 {
		return fmt.Errorf("error unlocking lock (%d): %v", res.StatusCode, err)
	}

	return nil
}
