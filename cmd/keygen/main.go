//go:build keygen
// +build keygen

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func generateAPIKey(secretPhrase string) string {
	h := hmac.New(sha256.New, []byte(secretPhrase))
	h.Write([]byte(secretPhrase))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
	} else {
		fmt.Println("Current working directory:", cwd)
	}

	// Load the .env file
	if err := godotenv.Load("configs/.env"); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Retrieve the secret phrase from the environment variable
	secretPhrase := os.Getenv("SECRET_PHRASE")
	if secretPhrase == "" {
		fmt.Println("SECRET_PHRASE environment variable is not set.")
		return
	}

	apiKey := generateAPIKey(secretPhrase)
	fmt.Println("Generated API Key:", apiKey)
}
