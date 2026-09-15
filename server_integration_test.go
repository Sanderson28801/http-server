package main

import (
	"io"
	"net/http"
	"testing"
	"time"
)

// TestThinVerticalSlice acts as a black-box client against your server.
// REQUIREMENT: You must build and run your HTTP server on port 8080 before running this test.
func TestThinVerticalSlice(t *testing.T) {
	// 1. Setup an HTTP client with a strict timeout to catch hanging connections
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	// 2. Construct the request
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/ping", nil)
	if err != nil {
		t.Fatalf("Test setup failed - could not create request: %v", err)
	}

	// 3. Execute the request against the black-box server
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to connect to server. Is it running on :8080? Error: %v", err)
	}
	// Ensure the TCP connection is closed after the test
	defer resp.Body.Close()

	// 4. Assert Protocol/Status Constraints
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
	}

	if resp.Proto != "HTTP/1.1" {
		t.Errorf("Expected HTTP/1.1 protocol, got %s", resp.Proto)
	}

	// 5. Assert Header Constraints
	if resp.Header.Get("Content-Length") != "4" {
		t.Errorf("Expected Content-Length header to be '4', got %q", resp.Header.Get("Content-Length"))
	}

	if resp.Header.Get("Content-Type") != "text/plain" {
		t.Errorf("Expected Content-Type to be 'text/plain', got %q", resp.Header.Get("Content-Type"))
	}

	// 6. Assert Body Constraints
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedBody := "pong"
	if string(bodyBytes) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(bodyBytes))
	}
}
