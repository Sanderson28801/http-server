package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// REQUIREMENT: Your server must be running on localhost:8080 before running these tests.

// Test 1: Validates that your parser extracts the Request Line and Headers,
// and dynamically injects them into the response body.
func TestParserAndEcho(t *testing.T) {
	client := http.Client{Timeout: 2 * time.Second}

	// Create a request with custom headers
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/echo-test", nil)
	if err != nil {
		t.Fatalf("Failed to build request: %v", err)
	}

	req.Header.Add("X-Custom-Header", "Architecture-Rocks")
	req.Header.Add("User-Agent", "Go-Test-Client")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Server connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	bodyStr := string(bodyBytes)

	// Assert the server successfully parsed and echoed the data
	expectedSubstrings := []string{
		"POST",               // The parsed method
		"/echo-test",         // The parsed URI
		"Architecture-Rocks", // The parsed custom header value
		"Go-Test-Client",     // The parsed User-Agent
	}

	for _, expected := range expectedSubstrings {
		if !strings.Contains(bodyStr, expected) {
			t.Errorf("Parser failed. Expected response body to contain %q, but it didn't.\nBody was:\n%s", expected, bodyStr)
		}
	}
}

// Test 2: Validates your server's ability to reject malformed HTTP requests gracefully.
func TestMalformedRequestRejection(t *testing.T) {
	// We use raw net.Dial instead of http.Client because http.Client won't let us send garbage data.
	conn, err := net.DialTimeout("tcp", "localhost:8080", 2*time.Second)
	if err != nil {
		t.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	// Send garbage data that ends with the valid HTTP terminator, but is not a valid HTTP request
	garbageRequest := "THIS IS NOT HTTP\r\nHost: localhost:8080\r\n\r\n"
	_, err = conn.Write([]byte(garbageRequest))
	if err != nil {
		t.Fatalf("Failed to write garbage data: %v", err)
	}

	// Read the response
	reader := bufio.NewReader(conn)
	statusLine, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read server response: %v", err)
	}

	// Assert the server recognized it as invalid and returned a 400 Bad Request
	if !strings.Contains(statusLine, "400 Bad Request") {
		t.Errorf("Expected status line to contain '400 Bad Request', got: %q", statusLine)
	}
}

// Test 3: Validates the concurrency model by firing 1,000 requests simultaneously.
func TestHighConcurrency(t *testing.T) {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1000,
		},
	}

	totalRequests := 1000
	var wg sync.WaitGroup
	wg.Add(totalRequests)

	errorsCount := 0
	var mu sync.Mutex

	for i := 0; i < totalRequests; i++ {
		go func(reqID int) {
			defer wg.Done()

			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/load-%d", reqID), nil)
			resp, err := client.Do(req)

			if err != nil {
				mu.Lock()
				errorsCount++
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				mu.Lock()
				errorsCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait() // Wait for all 1000 goroutines to finish

	if errorsCount > 0 {
		t.Fatalf("Concurrency test failed: %d out of %d requests resulted in errors or non-200 statuses. Ensure you are using goroutines for handleConnection.", errorsCount, totalRequests)
	}
}
