package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/aube/url-shortener/internal/app/router"
	"github.com/aube/url-shortener/internal/app/store"
)

// Run initializes and starts the URL shortener application.
// It performs the following steps:
//  1. Loads configuration using config.NewConfig()
//  2. Initializes the storage backend using store.NewStore()
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
	storage := store.NewStore()

	// Create router with all endpoints and middleware
	routes := router.New(storage, config.BaseURL)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	var server = http.Server{Addr: ":8080"}

	go func() {
		var err error

		if config.EnableHTTPS {

			log.Println("HTTP server starting", "address", config.ServerAddress)
			// Start HTTPS server
			err = http.ListenAndServeTLS(
				config.ServerAddress,
				config.PublicCertFile,
				config.PrivateCertFile,
				routes)
		} else {

			log.Println("HTTPS server starting", "address", config.ServerAddress)
			// Start HTTP server
			err = http.ListenAndServe(config.ServerAddress, routes)
		}

		if err != nil {
			log.Fatal("Starting server", "err", err)
		}
	}()

	sig := <-sigs
	log.Printf("Received signal: %v\n", sig)

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Close database
	if err := store.Close(); err != nil {
		log.Printf("Store close error: %v", err)
		return err
	} else {
		log.Printf("Store connection closed")
	}

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v\n", err)
		return err
	} else {
		log.Println("Server gracefully stopped")
	}
	return nil
}
