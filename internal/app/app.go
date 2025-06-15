package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/aube/url-shortener/internal/app/router"
	"github.com/aube/url-shortener/internal/app/store"
	"github.com/go-chi/chi/v5"
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
	r := router.New(storage, config.BaseURL)

	err := startServer(config, &r)

	if err != nil {
		log.Fatal("Starting server", "err", err)
	}

	return nil
}

func startServer(config config.EnvConfig, r *chi.Router) error {
	var err error

	// Construct server address from config
	log.Println("Server starting", "address", config.ServerAddress)

	var srv = http.Server{Addr: ":8080"}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	stopServer(&srv, sigs)

	fmt.Printf("TYPE: %T\n", &srv) // "int"

	if config.EnableHTTPS {
		// Start HTTPS server
		err = http.ListenAndServeTLS(
			config.ServerAddress,
			config.PublicCertFile,
			config.PrivateCertFile,
			*r)
	} else {
		// Start HTTP server
		err = http.ListenAndServe(config.ServerAddress, *r)
	}

	return err
}

func stopServer(srv *http.Server, sigs chan os.Signal) {
	for {
		time.Sleep(1 * time.Second)
		if sig, ok := <-sigs; ok {
			if err := store.Close(); err != nil {
				log.Printf("Store close error: %v", err)
			}
			if err := srv.Shutdown(context.Background()); err != nil {
				log.Printf("Server shutdown error: %v", err)
			}
			log.Println("Server stopped.", "Signal:", sig)
			return
		}
	}
}
