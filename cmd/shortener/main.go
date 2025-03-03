package main

import (
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/cmd/shortener/config"
	"github.com/aube/url-shortener/internal/app/handlers"
	"github.com/aube/url-shortener/internal/app/middlewares"
	"github.com/aube/url-shortener/internal/app/store"
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

	// r.Post("/*", middlewares.LoggingMiddleware(middlewares.GzipMiddleware(handlers.HandlerRoot(config.BaseURL))))
	// r.Post("/api/*", middlewares.LoggingMiddleware(middlewares.GzipMiddleware(handlers.HandlerAPI(config.BaseURL))))
	// r.Get("/api/user/urls", middlewares.LoggingMiddleware(handlers.HandlerAPIUserUrls(config.BaseURL)))
	// r.Get("/{id}", middlewares.LoggingMiddleware(middlewares.GzipMiddleware(handlers.HandlerID())))

	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.LoggingMiddleware,
			middlewares.GzipMiddleware,
		)
		r.Get("/{id}", handlers.HandlerID())
		r.Post("/*", handlers.HandlerRoot(config.BaseURL))
		r.Post("/api/*", handlers.HandlerAPI(config.BaseURL))
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.LoggingMiddleware)
		r.Get("/api/user/urls", handlers.HandlerAPIUserUrls(config.BaseURL))
	})

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
