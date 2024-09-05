package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestBuildRequest_Success tests BuildRequest for successful request creation.
func TestBuildRequest_Success(t *testing.T) {
	method := "POST"
	url := "https://example.com/api"
	body := []byte(`{"key": "value"}`)

	req, err := BuildRequest(method, url, body)
	if err != nil {
		t.Fatalf("BuildRequest returned an error: %v", err)
	}

	if req.Method != method {
		t.Errorf("Expected method %s, got %s", method, req.Method)
	}
	if req.URL.String() != url {
		t.Errorf("Expected URL %s, got %s", url, req.URL.String())
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(req.Body)
	if err != nil {
		t.Fatalf("Error")
	}
	if buf.String() != string(body) {
		t.Errorf("Expected body %s, got %s", body, buf.String())
	}
}

// TestBuildRequest_Error tests BuildRequest with invalid URL.
func TestBuildRequest_Error(t *testing.T) {
	invalidURL := "://invalid-url"
	_, err := BuildRequest("GET", invalidURL, nil)
	if err == nil {
		t.Fatalf("BuildRequest did not return an error for an invalid URL")
	}
}

// TestSendRequest_Success tests SendRequest for a successful request and response.
func TestSendRequest_Success(t *testing.T) {
	// Create a test server with a simple handler
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"result": "success"}`))
		if err != nil {
			t.Fatalf("Test Error")
		}
	}))
	defer ts.Close()

	// Create a request to the test server
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Test SendRequest
	body, status, err := SendRequest(req)
	if err != nil {
		t.Fatalf("SendRequest returned an error: %v", err)
	}
	if status != "200 OK" {
		t.Errorf("Expected status %s, got %s", "200 OK", status)
	}
	if body != `{"result": "success"}` {
		t.Errorf("Expected body %s, got %s", `{"result": "success"}`, body)
	}
}

// TestSendRequest_Error tests SendRequest for an error case.
func TestSendRequest_Error(t *testing.T) {
	// Create an invalid request
	req, err := http.NewRequest("GET", "http://invalid-url", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Test SendRequest
	_, _, err = SendRequest(req)
	if err == nil {
		t.Fatalf("SendRequest did not return an error for an invalid URL")
	}
}

// TestAddHeaders tests AddHeaders for correctly adding headers to a request.
func TestAddHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)

	headers := map[string]string{
		"Authorization": "Bearer token",
		"User-Agent":    "my-app",
	}

	AddHeaders(req, headers)

	// Check if headers are correctly added
	for key, value := range headers {
		if got := req.Header.Get(key); got != value {
			t.Errorf("Expected header %s to have value %s, got %s", key, value, got)
		}
	}
}
