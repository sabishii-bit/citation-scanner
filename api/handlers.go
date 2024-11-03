package api

import (
	"citation-scanner/internal/cache"
	"citation-scanner/internal/parser"
	"encoding/json"
	"net/http"
)

// Home handler just for the base route
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Successfully touched the API."))
}

func parsePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the incoming JSON payload
	var requestBody struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if URL is provided
	if requestBody.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Check cache for the URL
	cachedResponse, found, err := cache.GetCachedResponse(requestBody.URL)
	if err != nil {
		http.Error(w, "Error checking cache: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if found {
		// Cached response found, return it
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedResponse))
		return
	}

	// No cached response or it's expired, call the ParsePageClaims function
	parsedClaims, err := parser.ParsePageClaims(requestBody.URL)
	if err != nil {
		http.Error(w, "Failed to parse page: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the parsed claims to JSON
	responseData, err := json.Marshal(parsedClaims)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Cache the new response
	err = cache.CacheResponse(requestBody.URL, string(responseData))
	if err != nil {
		http.Error(w, "Failed to cache response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the parsed claims in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}
