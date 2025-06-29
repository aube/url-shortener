package app

import (
	"log"
	"net/http"

	restapi "github.com/aube/url-shortener/internal/api/rest"
	"github.com/aube/url-shortener/internal/app/config"
)

func httpServerStart() *http.Server {
	// Load application configuration
	config := config.NewConfig()

	// Create router with all endpoints and middleware
	routes := restapi.New(config.BaseURL)

	var server = http.Server{Addr: ":8080"}
	var err error

	if config.EnableHTTPS {
		log.Println("HTTPS server starting", "address", config.ServerAddress)
		// Start HTTPS server
		err = http.ListenAndServeTLS(
			config.ServerAddress,
			config.PublicCertFile,
			config.PrivateCertFile,
			routes)
	} else {

		log.Println("HTTP server starting", "address", config.ServerAddress)
		// Start HTTP server
		err = http.ListenAndServe(config.ServerAddress, routes)
	}

	if err != nil {
		log.Fatal("Starting server", "err", err)
	}

	return &server
}
