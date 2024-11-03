package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const port = ":4145"

// StartServer starts the HTTP API server
func StartServer() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(APIKeyMiddleware) // Use API key middleware for authentication; applies to all routes

	// Setup routes from routes.go
	routes(r)

	fmt.Println("Starting server on port", port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
