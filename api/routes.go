package api

import (
	"github.com/go-chi/chi/v5"
)

// routes sets up all the API routes
func routes(r chi.Router) {
	r.Get("/", homeHandler)
	r.Post("/parse", parsePageHandler)
}
