package parser

import (
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
