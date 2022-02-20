package reporter

import (
	"bytes"
	"challenge/go-healthcheck/checker"
	"challenge/go-healthcheck/client"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	host   string = "https://backend-challenge.line-apps.com"
	method string = "POST"
)

// Create report request
func CreateRequest(accessToken string, pingResult checker.Result) (*http.Request, error) {
	json_data, _ := json.Marshal(pingResult)

	url := fmt.Sprintf("%s/healthcheck/report", host)
	request, newReqErr := http.NewRequest(method, url, bytes.NewBuffer(json_data))
	if newReqErr != nil {
		return nil, newReqErr
	}

	authorizationHeader := fmt.Sprintf("Bearer %s", accessToken)

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Authorization", authorizationHeader)

	return request, nil
}

// Report ping result to Healthcheck API
func Report(client client.HttpClient, accessToken string, pingResult checker.Result) error {
	request, createReqErr := CreateRequest(accessToken, pingResult)
	if createReqErr != nil {
		return createReqErr
	}

	resp, reqErr := client.Do(request)
	if reqErr != nil {
		return reqErr
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode >= 400 {
		return errors.New("Failed to report healthcheck status")
	}

	return nil
}
