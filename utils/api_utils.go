package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// BuildRequest creates an HTTP request with the given method, URL, and body.
// Returns the request and an error if creation fails.
func BuildRequest(method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	return req, nil
}

// SendRequest sends an HTTP request and returns the response body, status, and any error.
// Returns the response body as a string, the HTTP status code as a string, and an error if the request fails.
func SendRequest(req *http.Request) (string, string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(responseBody), resp.Status, nil
}

// AddHeaders sets multiple headers on an HTTP request.
// It takes a map of header key-value pairs and applies them to the request.
func AddHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}
