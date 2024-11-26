package parser

import (
	"citation-scanner/internal/cache"
	"encoding/json"
	"fmt"
	"testing"
)

// TestParsePageClaims tests the ParsePageClaims function by scraping a real page and using the OpenAI client.
func TestParsePageClaims(t *testing.T) {
	// Define the page URL to scrape and parse
	url := "https://en.wikipedia.org/wiki/Go_(programming_language)"

	// Parse the claims using the parser package
	parsedClaims, err := ParsePageClaims(url)
	if err != nil {
		t.Fatalf("Error parsing claims: %v", err)
	}

	// Print the parsed claims as JSON
	claimsJSON, err := json.MarshalIndent(parsedClaims, "", "  ")
	if err != nil {
		t.Fatalf("Error marshalling claims to JSON: %v", err)
	}
	fmt.Println(string(claimsJSON))
}

// TestParseAndAggregateClaims tests the ParseAndAggregateClaims function by parsing a root page and its sources.
func TestParseAndAggregateClaims(t *testing.T) {
	// Initialize cache for testing
	err := cache.InitializeCache()
	if err != nil {
		t.Fatalf("Failed to initialize cache: %v", err)
	}
	defer cache.CloseCache()

	// Define the root page URL for the test
	rootURL := "https://en.wikipedia.org/wiki/High-trust_and_low-trust_societies"

	// Define the maximum depth for recursion
	maxDepth := 1

	// Call the ParseAndAggregateClaims function
	aggregatedClaims, err := ParseAndAggregateClaims(rootURL, maxDepth)
	if err != nil {
		t.Fatalf("Error parsing and aggregating claims: %v", err)
	}

	// Print the aggregated claims as JSON
	claimsJSON, err := json.MarshalIndent(aggregatedClaims, "", "  ")
	if err != nil {
		t.Fatalf("Error marshalling aggregated claims to JSON: %v", err)
	}
	fmt.Println(string(claimsJSON))

	// Validate the aggregated claims
	if len(aggregatedClaims.AllClaims) == 0 {
		t.Error("Expected to find claims, but none were found")
	}

	// Check for errors in the aggregated results
	if len(aggregatedClaims.Errors) > 0 {
		t.Logf("Encountered errors during parsing: %v", aggregatedClaims.Errors)
	}
}
