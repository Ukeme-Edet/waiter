package main

import (
	"fmt"
	"os"
	"waiter/server"
)

var dir = ""

func main() {
	// Create a new server instance
	srv := &server.Server{}

	for i, arg := range os.Args {
		if arg == "--directory" {
			if i+1 >= len(os.Args) {
				fmt.Println("Invalid command-line arguments")
				os.Exit(2)
			}
			dir = os.Args[i+1]
		}
	}

	// Start the server on port 4221
	if err := srv.Run("0.0.0.0:4221", dir); err != nil {
		fmt.Println("Error running server:", err)
		os.Exit(1)
	}
	fmt.Println("Server stopped gracefully")
	os.Exit(0)
}
