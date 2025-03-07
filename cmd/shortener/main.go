package main

import (
	"net/http"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/aube/url-shortener/internal/app/handlers"
	"github.com/aube/url-shortener/internal/app/middlewares"
	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
)

func main() {
	config := config.NewConfig()

	var storage store.Storage
	if config.FileStoragePath != "" {
		storage = store.NewFileStore(config.FileStoragePath)
	} else {
		storage = store.NewMemStore()
	}

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.LoggingMiddleware,
			middlewares.GzipMiddleware,
		)
		r.Get("/{id}", handlers.HandlerID(storage))
		r.Post("/*", handlers.HandlerRoot(storage, config.BaseURL))
		r.Post("/api/*", handlers.HandlerAPI(storage, config.BaseURL))
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.LoggingMiddleware)
		r.Get("/api/user/urls", handlers.HandlerAPIUserUrls(storage, config.BaseURL))
	})

	// empty handler for prevent error on automatic browser favicon request
	r.Get("/favicon.ico", http.HandlerFunc(handlers.HandlerEmpty))

	if err := logger.Initialize(); err != nil {
		return
	}

	address := config.ServerHost + ":" + config.ServerPort
	logger.Infoln("Server starting at", address)

	err := http.ListenAndServe(address, r)

	if err != nil {
		logger.Infoln("Error starting server:", err)
	}

}
