package webscraper

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// ScrapePage fetches the content of a webpage and extracts the text content from its HTML.
func ScrapePage(url string) (string, error) {
	// Make a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	// Parse the HTML document
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse the page HTML: %v", err)
	}

	// Extract text content from the parsed HTML
	var sb strings.Builder
	scrapeText(doc, &sb)

	return sb.String(), nil
}

// scrapeText extracts text nodes recursively from an HTML node and writes them to the StringBuilder.
func scrapeText(n *html.Node, sb *strings.Builder) {
	if n.Type == html.TextNode {
		sb.WriteString(n.Data)
		sb.WriteString(" ")
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		scrapeText(c, sb)
	}
}

// ScrapeBody fetches only the text content within the <body> tag of a webpage.
func ScrapeBody(url string) (string, error) {
	// Make a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	// Parse the HTML document
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse the page HTML: %v", err)
	}

	// Find the <body> tag
	bodyNode := findBodyNode(doc)
	if bodyNode == nil {
		return "", fmt.Errorf("no <body> tag found in the page")
	}

	// Extract text content only within the <body> tag
	var sb strings.Builder
	scrapeText(bodyNode, &sb)

	return sb.String(), nil
}

// findBodyNode finds and returns the <body> node in an HTML document.
func findBodyNode(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "body" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if body := findBodyNode(c); body != nil {
			return body
		}
	}
	return nil
}
