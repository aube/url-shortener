package main

import (
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/cmd/shortener/config"
	"github.com/aube/url-shortener/internal/app/handlers"
	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/gzip"
	"github.com/aube/url-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
)

func init() {
	config := config.NewConfig()
	store.NewFileStore(config.FileStoragePath)
}

func main() {
	config := config.NewConfig()
	r := chi.NewRouter()

	r.Post("/*", logger.LoggingMiddleware(gzip.GzipMiddleware(handlers.HandlerRoot(config.BaseURL))))
	r.Post("/api/*", logger.LoggingMiddleware(gzip.GzipMiddleware(handlers.HandlerAPI(config.BaseURL))))
	r.Get("/api/user/urls", logger.LoggingMiddleware(handlers.HandlerAPIUserUrls(config.BaseURL)))
	r.Get("/{id}", logger.LoggingMiddleware(gzip.GzipMiddleware(handlers.HandlerID())))

	// empty handler for prevent error on automatic browser favicon request
	r.Get("/favicon.ico", http.HandlerFunc(handlers.HandlerEmpty))

	if err := logger.Initialize(); err != nil {
		return
	}

	err := http.ListenAndServe(config.ServerHost+":"+config.ServerPort, r)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
