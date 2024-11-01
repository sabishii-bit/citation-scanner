package main

import (
	"citation-scanner/api"
	"fmt"
)

func main() {
	fmt.Println("-- Program Entry --")

	// Start the HTTP server
	api.StartServer()

	fmt.Println("-- Program Exit --")
}
