package main

import (
	"citation-scanner/api"
	"citation-scanner/internal/cache"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("-- Program Entry --")

	// Initialize the cache database
	if err := cache.InitializeCache(); err != nil {
		fmt.Printf("Failed to initialize cache: %v\n", err)
		return
	}
	defer cache.CloseCache()

	// Set up a channel to listen for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the HTTP server in a goroutine
	go func() {
		api.StartServer()
	}()

	// Wait for a termination signal
	<-quit
	fmt.Println("Shutting down gracefully...")

	// Program will exit here, and deferred CloseCache will be called
	fmt.Println("-- Program Exit --")
}
