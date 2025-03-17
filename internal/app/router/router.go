package router

import (
	"context"
	"net/http"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/aube/url-shortener/internal/app/handlers"
	"github.com/aube/url-shortener/internal/app/middlewares"
	"github.com/go-chi/chi/v5"
)

type StorageGet interface {
	Get(ctx context.Context, key string) (value string, ok bool)
}
type StorageList interface {
	List(ctx context.Context) map[string]string
}
type StoragePing interface {
	Ping() error
}
type StorageSet interface {
	Set(ctx context.Context, key string, value string) error
}
type StorageSetMultiple interface {
	SetMultiple(ctx context.Context, l map[string]string) error
}

type Storage interface {
	StorageGet
	StorageList
	StoragePing
	StorageSet
	StorageSetMultiple
}

func Connect(ctx context.Context, storage Storage) chi.Router {
	config := config.NewConfig()
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

	return r
}
