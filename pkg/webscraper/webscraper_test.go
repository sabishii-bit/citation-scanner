package webscraper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestScrapePage is a test function for ScrapePage.
func TestScrapePage(t *testing.T) {
	// Create a mock server to serve HTML content for testing
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<html>
				<head><title>Test Page</title></head>
				<body>
					<p>Hello, World!</p>
					<p>This is a test page.</p>
				</body>
			</html>
		`))
	}))
	defer mockServer.Close()

	// Use the mock server's URL for testing
	result, err := ScrapePage(mockServer.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that the result contains the expected text content
	expectedText := "Hello, World! This is a test page."
	result = strings.TrimSpace(result)

	if !strings.Contains(result, "Hello, World!") || !strings.Contains(result, "This is a test page.") {
		t.Errorf("expected '%s' to contain '%s'", result, expectedText)
	}
}

// TestScrapePageErrorHandling is a test function for error handling in ScrapePage.
func TestScrapePageErrorHandling(t *testing.T) {
	// Test with an invalid URL
	_, err := ScrapePage("http://invalid-url")
	if err == nil {
		t.Errorf("expected error for invalid URL, got nil")
	}

	// Test with a non-200 HTTP status code
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	_, err = ScrapePage(mockServer.URL)
	if err == nil {
		t.Errorf("expected error for non-200 status code, got nil")
	}
}

// TestScrapeRealPage is a test function that scrapes a real Wikipedia page.
func TestScrapeRealPage(t *testing.T) {
	// URL to scrape (a real Wikipedia page)
	url := "https://en.wikipedia.org/wiki/Go_(programming_language)"

	// Scrape the Wikipedia page
	result, err := ScrapePage(url)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Print the result
	fmt.Println("Scraped Content:")
	fmt.Println(result)

	// Assume a 200 response is a success
	if result == "" {
		t.Errorf("expected non-empty result, got empty string")
	}
}
