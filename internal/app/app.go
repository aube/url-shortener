package app

import (
	"log"
	"net/http"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/aube/url-shortener/internal/app/router"
	"github.com/aube/url-shortener/internal/app/store"
)

// Run initializes and starts the URL shortener application.
// It performs the following steps:
//  1. Loads configuration using config.NewConfig()
//  2. Initializes the storage backend using store.MewStore()
//  3. Creates the router with all endpoints using router.New()
//  4. Starts the HTTP server on the configured address
//
// The function blocks until the server exits and returns any error that occurs.
// In case of a fatal server error, it logs the error and terminates the program.
//
// Example usage:
//
//	if err := app.Run(); err != nil {
//	    log.Fatal("Application failed:", err)
//	}
func Run() error {
	// Load application configuration
	config := config.NewConfig()

	// Initialize storage (database, file, or memory based on config)
	storage := store.MewStore()

	// Create router with all endpoints and middleware
	r := router.New(storage, config.BaseURL)

	// Construct server address from config
	address := config.ServerHost + ":" + config.ServerPort
	log.Println("Server starting", "address", address)

	// Start HTTP server
	err := http.ListenAndServe(address, r)

	if err != nil {
		log.Fatal("Starting server", "err", err)
	}

	return nil
}
