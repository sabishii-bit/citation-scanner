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

// TestCacheResponseAndGetCachedResponse tests storing and retrieving a cached response and checks if TTL is respected.
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
	err = CacheResponse(url, response)
	if err != nil {
		t.Fatalf("Failed to cache response: %v", err)
	}

	// Retrieve the cached response immediately, which should not be expired
	cachedResponse, found, err := GetCachedResponse(url)
	if err != nil {
		t.Fatalf("Failed to get cached response: %v", err)
	}
	if !found {
		t.Fatalf("Expected cache entry for URL %s, but it was not found", url)
	}
	if cachedResponse != response {
		t.Errorf("Expected cached response '%s', got '%s'", response, cachedResponse)
	}

	// Wait for the TTL to expire
	time.Sleep(CacheTTL + 1*time.Second)

	// Attempt to retrieve the response after TTL, should return as expired
	cachedResponse, found, err = GetCachedResponse(url)
	if err != nil {
		t.Fatalf("Error retrieving cached response: %v", err)
	}
	if found {
		t.Errorf("Expected cache entry to be expired for URL %s, but it was found", url)
	}
}

// TestCacheUpdate verifies that caching a response updates the timestamp and respects TTL.
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
	err = CacheResponse(url, initialResponse)
	if err != nil {
		t.Fatalf("Failed to cache initial response: %v", err)
	}

	// Retrieve the initial cached response immediately
	cachedResponse, found, err := GetCachedResponse(url)
	if err != nil || !found {
		t.Fatalf("Failed to retrieve initial cached response")
	}
	if cachedResponse != initialResponse {
		t.Errorf("Expected cached response '%s', got '%s'", initialResponse, cachedResponse)
	}

	// Wait for a moment to differentiate timestamps
	time.Sleep(1 * time.Second)

	// Store the updated response, which should refresh the timestamp
	err = CacheResponse(url, updatedResponse)
	if err != nil {
		t.Fatalf("Failed to cache updated response: %v", err)
	}

	// Retrieve the updated cached response
	cachedResponse, found, err = GetCachedResponse(url)
	if err != nil || !found {
		t.Fatalf("Failed to retrieve updated cached response")
	}
	if cachedResponse != updatedResponse {
		t.Errorf("Expected updated cached response '%s', got '%s'", updatedResponse, cachedResponse)
	}

	// Verify that the updated entry is still valid within the TTL
	time.Sleep(CacheTTL / 2)
	cachedResponse, found, err = GetCachedResponse(url)
	if err != nil || !found {
		t.Fatalf("Expected cache entry to be valid within TTL, but it was expired or not found")
	}
	if cachedResponse != updatedResponse {
		t.Errorf("Expected cached response '%s', got '%s'", updatedResponse, cachedResponse)
	}
}
