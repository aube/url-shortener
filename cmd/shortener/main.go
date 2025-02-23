package main

import (
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/internal/app/handlers"
	"github.com/aube/url-shortener/internal/gzip"
	"github.com/aube/url-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
)

func main() {
	config := NewConfig()

	r := chi.NewRouter()

	r.Post("/*", logger.LoggingMiddleware(gzip.GzipMiddleware(handlers.HandlerRoot(config.BaseURL))))
	r.Get("/{id}", logger.LoggingMiddleware(gzip.GzipMiddleware(handlers.HandlerID())))

	// empty handler for prevent error on automatic browser favicon request
	r.Get("/favicon.ico", http.HandlerFunc(handlers.HandlerEmpty))

	if err := logger.Initialize("debug"); err != nil {
		return
	}

	logger.Log.Info("Running server", zap.String("address", config.ServerAddress))

	err := http.ListenAndServe(config.ServerAddress+":"+config.ServerPort, r)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
