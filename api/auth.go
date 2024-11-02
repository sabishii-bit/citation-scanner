package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Load the environment variables from the .env file
	if err := godotenv.Load("configs/.env"); err != nil {
		panic("Error loading .env file")
	}
}

// APIKeyMiddleware checks for a valid API key in the request headers.
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		secretPhrase := os.Getenv("SECRET_PHRASE")

		if secretPhrase == "" {
			http.Error(w, "Server misconfiguration", http.StatusInternalServerError)
			return
		}

		// Validate the API key using HMAC with the secret phrase
		if validateHMAC(apiKey, secretPhrase) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}

// validateHMAC validates the API key using HMAC with a secret phrase.
func validateHMAC(apiKey, secretPhrase string) bool {
	h := hmac.New(sha256.New, []byte(secretPhrase))
	h.Write([]byte(secretPhrase))
	expectedMAC := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(apiKey), []byte(expectedMAC))
}
