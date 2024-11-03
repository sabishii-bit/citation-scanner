package api

import (
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

	// Call the ParsePageClaims function
	parsedClaims, err := parser.ParsePageClaims(requestBody.URL)
	if err != nil {
		http.Error(w, "Failed to parse page: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the parsed claims in JSON format
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(parsedClaims); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
