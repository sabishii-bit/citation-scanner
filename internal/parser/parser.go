package parser

import (
	"citation-scanner/pkg/openai"
	"citation-scanner/pkg/webscraper"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// ParsedClaims represents the structure of the JSON object for claims and sources.
type ParsedClaims struct {
	Page   string  `json:"page"`
	Claims []Claim `json:"claims"`
}

// Claim represents a single claim and its source.
type Claim struct {
	Claim  string `json:"claim"`
	Source string `json:"source"`
}

// ParsePageClaims takes a URL, scrapes the content, and uses OpenAI to extract claims and their sources.
func ParsePageClaims(url string) (*ParsedClaims, error) {
	// Load the .env file located under ./configs/.env
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("Error getting current working directory: %v", err)
	}
	envPath := filepath.Join(cwd, "..", "..", "configs", ".env")
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	// Get the OpenAI API key from environment variables
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Environment variable OPENAI_API_KEY is required but not set")
	}

	// Create a new OpenAIClient instance with customized settings
	openAIClient := openai.NewClient(apiKey,
		openai.WithTemperature(0.15),
		openai.WithSystemRole("You are an expert in extracting claims from articles."),
	)

	// Step 1: Scrape the content of the page using the webscraper package
	scrapedContent, err := webscraper.ScrapePage(url)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape the page: %v", err)
	}

	// Step 2: Prepare the prompt for OpenAI to identify claims and their sources
	prompt := fmt.Sprintf(`
		You are a parser that extracts claims and their sources from a scraped article.
		Please read the following content and provide all the claims and their corresponding sources linked from the page.
		Make sure that the claims extracted are direct quotes from the scraped page text.
		Provide the actual citation links to the relevant sources, not the reference markers.
		DO NOT wrap response with Markdown code-block formatting.
		ONLY respond with the JSON object, do not include any additional text.
		Content: "%s"
		Example response:
		{
			"claims": [
				{"claim": "... Example claim 1", "source": "https://www.example-source-1.com/"},
				{"claim": "... Example claim 2 ...", "source": "https://www.example-source-2.com/"}
			]
		}
	`, scrapedContent)

	// Step 3: Use OpenAIClient to get claims from the scraped content
	response, err := openAIClient.SendChatRequest(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to extract claims using OpenAI: %v", err)
	}

	fmt.Println("Parsed JSON:")
	fmt.Println(response)

	// Step 4: Parse the JSON response
	var parsedClaims ParsedClaims
	err = json.Unmarshal([]byte(response), &parsedClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response as JSON: %v", err)
	}

	// Step 5: Set the page URL in the parsed claims
	parsedClaims.Page = url

	return &parsedClaims, nil
}
