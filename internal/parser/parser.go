package parser

import (
	"citation-scanner/pkg/openai"
	"citation-scanner/pkg/webscraper"
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// ParsedClaims represents the structure of the JSON object for claims and sources.
type ParsedClaims struct {
	Page   string  `json:"page"`
	Claims []Claim `json:"claims"`
}

// Claim represents a single claim and its source.
type Claim struct {
	Claim  string   `json:"claim"`
	Source []string `json:"sources"`
}

// ParsePageClaims takes a URL, scrapes the content, and uses OpenAI to extract claims and their sources.
func ParsePageClaims(url string) (*ParsedClaims, error) {

	// Load the .env file
	if err := godotenv.Load("configs/.env"); err != nil {
		fmt.Println("Error loading .env file")
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	// Get the OpenAI API key from environment variables
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Environment variable OPENAI_API_KEY is required but not set")
	}

	// Create a new OpenAIClient instance with customized settings
	openAIClient := openai.NewClient(apiKey,
		openai.WithTemperature(0.1),
		openai.WithSystemRole("You are an expert in extracting claims from articles."),
	)

	// Step 1: Scrape the content of the page using the webscraper package
	scrapedContent, err := webscraper.ScrapeBody(url)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape the page: %v", err)
	}

	// Step 2: Prepare the prompt for OpenAI to identify claims and their sources
	prompt := fmt.Sprintf(`
		You are a parser that extracts claims and their reference sources from a scraped webpage article.
		Please read the following content and provide ALL of the claims, and their corresponding sources linked from the page.
		Sources are identified by <a> tags in a claim or reference marker(s). All sources must be returned and associated to a claim. 
		There can be more than one source to a claim, so return them in an array of strings.
		Make sure that the claims extracted are direct quotes from the scraped page text; prefix and/or postfix with "..." if a quoted claim is a section of a sentence.
		Provide the actual citation links to the associated sources, not the reference markers.
		DO NOT wrap response with Markdown code-block formatting. DO NOT omit any claims or sources from the content in your response.
		ALL CLAIMS AND SOURCES MUST BE RETURNED, REGARDLESS OF PROCESSING TIME OR LENGTH OF RESPONSE.
		Respond only with a JSON object formatted as follows:
		{
			"claims": [
				{"claim": "... Example claim 1[34][35].", "sources": ["https://www.example-source-1.com/article1", "https://www.example-source-1.org/"]},
				{"claim": "... Example claim 2[65] ...", "sources": ["https://www.example-source-2.com/"]}
			]
		}
		Content: "%s"
	`, scrapedContent)

	// Step 3: Use OpenAIClient to get claims from the scraped content
	response, err := openAIClient.SendChatRequest(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to extract claims using OpenAI: %v", err)
	}

	fmt.Println(response)

	// Step 4: Parse the JSON response
	var parsedClaims ParsedClaims
	err = json.Unmarshal([]byte(response), &parsedClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response as JSON: %v", err)
	}

	// Ensure that each claim's Source is an empty array if it's nil
	for i := range parsedClaims.Claims {
		if parsedClaims.Claims[i].Source == nil {
			parsedClaims.Claims[i].Source = []string{}
		}
	}

	// Step 5: Set the page URL in the parsed claims
	parsedClaims.Page = url

	return &parsedClaims, nil
}
