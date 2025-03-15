package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/aube/url-shortener/internal/app/handlers"
	"github.com/aube/url-shortener/internal/app/middlewares"
	"github.com/aube/url-shortener/internal/app/store"
	"github.com/aube/url-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
)

type Storage interface {
	Get(ctx context.Context, key string) (value string, ok bool)
	List(ctx context.Context) map[string]string
	Ping() error
	Set(ctx context.Context, key string, value string) error
	SetMultiple(ctx context.Context, l map[string]string) error
}

func Run() error {
	config := config.NewConfig()

	var storage Storage

	if config.DatabaseDSN != "" {
		storage = store.NewDBStore(config.DatabaseDSN)
	} else if config.FileStoragePath != "" {
		storage = store.NewFileStore(config.FileStoragePath)
	} else {
		storage = store.NewMemStore()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.LoggingMiddleware,
			middlewares.GzipMiddleware,
		)
		r.Get("/{id}", handlers.HandlerID(ctx, storage))
		r.Post("/*", handlers.HandlerRoot(ctx, storage, config.BaseURL))
		r.Post("/api/*", handlers.HandlerAPI(ctx, storage, config.BaseURL))
		r.Post("/api/shorten/batch", handlers.HandlerShortenBatch(ctx, storage, config.BaseURL))
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.LoggingMiddleware)
		r.Get("/api/user/urls", handlers.HandlerAPIUserUrls(ctx, storage, config.BaseURL))
		r.Get("/ping", handlers.HandlerPing(storage))
	})

	// empty handler for prevent error on automatic browser favicon request
	r.Get("/favicon.ico", http.HandlerFunc(handlers.HandlerEmpty))

	if err := logger.Initialize(); err != nil {
		return errors.New("logger start error")
	}

	address := config.ServerHost + ":" + config.ServerPort
	logger.Infoln("Server starting at", address)

	err := http.ListenAndServe(address, r)

	if err != nil {
		logger.Infoln("Error starting server:", err)
		return err
	}

	return nil
}
