package cache

import (
	"testing"
	"time"
)

// TestInitializeCache verifies that the cache can be initialized and the table is created.
func TestInitializeCache(t *testing.T) {
	err := InitializeCache()
	if err != nil {
		t.Fatalf("Failed to initialize cache: %v", err)
	}
	defer CloseCache() // Ensure database is closed after test
}

// TestCacheResponseAndGetCachedResponse tests storing and retrieving a cached response.
func TestCacheResponseAndGetCachedResponse(t *testing.T) {
	err := InitializeCache()
	if err != nil {
		t.Fatalf("Failed to initialize cache: %v", err)
	}
	defer CloseCache()

	// Define test URL and response data
	url := "https://example.com"
	response := "{\"message\": \"test response\"}"

	// Store the response in the cache
	err = cacheResponse(url, response)
	if err != nil {
		t.Fatalf("Failed to cache response: %v", err)
	}

	// Retrieve the cached response
	cachedResponse, timestamp, found, err := getCachedResponse(url)
	if err != nil {
		t.Fatalf("Failed to get cached response: %v", err)
	}
	if !found {
		t.Fatalf("Expected cache entry for URL %s, but it was not found", url)
	}
	if cachedResponse != response {
		t.Errorf("Expected cached response '%s', got '%s'", response, cachedResponse)
	}

	// Check if the timestamp is recent
	parsedTimestamp, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Fatalf("Failed to parse timestamp: %v", err)
	}
	if time.Since(parsedTimestamp) > time.Minute {
		t.Errorf("Timestamp is too old, expected recent timestamp but got %s", timestamp)
	}
}

// TestCacheUpdate verifies that caching a response updates the timestamp.
func TestCacheUpdate(t *testing.T) {
	err := InitializeCache()
	if err != nil {
		t.Fatalf("Failed to initialize cache: %v", err)
	}
	defer CloseCache()

	// Define test URL and initial response data
	url := "https://example.com"
	initialResponse := "{\"message\": \"initial response\"}"
	updatedResponse := "{\"message\": \"updated response\"}"

	// Store the initial response
	err = cacheResponse(url, initialResponse)
	if err != nil {
		t.Fatalf("Failed to cache initial response: %v", err)
	}

	// Retrieve the initial timestamp
	_, initialTimestamp, found, err := getCachedResponse(url)
	if err != nil || !found {
		t.Fatalf("Failed to retrieve initial cached response")
	}

	// Wait a moment to ensure the timestamp difference
	time.Sleep(1 * time.Second)

	// Store the updated response
	err = cacheResponse(url, updatedResponse)
	if err != nil {
		t.Fatalf("Failed to cache updated response: %v", err)
	}

	// Retrieve the updated timestamp and response
	cachedResponse, updatedTimestamp, found, err := getCachedResponse(url)
	if err != nil || !found {
		t.Fatalf("Failed to retrieve updated cached response")
	}
	if cachedResponse != updatedResponse {
		t.Errorf("Expected updated cached response '%s', got '%s'", updatedResponse, cachedResponse)
	}

	// Verify the timestamp has changed
	if initialTimestamp == updatedTimestamp {
		t.Errorf("Expected timestamp to be updated, but it did not change")
	}
}
