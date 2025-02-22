package main

import (
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/internal/app/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	config := NewConfig()

	r := chi.NewRouter()

	r.Post("/*", handlers.HandlerRoot(config.BaseURL))
	r.Get("/{id}", handlers.HandlerID())

	// empty handler for prevent error on automatic browser favicon request
	r.Get("/favicon.ico", http.HandlerFunc(handlers.HandlerEmpty))

	if err := logger.Initialize("info"); err != nil {
		return
	}

	logger.Log.Info("Running server", zap.String("address", config.ServerAddress))

	err := http.ListenAndServe(config.ServerAddress+":"+config.ServerPort, r)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
