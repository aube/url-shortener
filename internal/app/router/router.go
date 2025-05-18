package router

import (
	"context"
	"net/http"

	"github.com/aube/url-shortener/internal/app/handlers"
	"github.com/aube/url-shortener/internal/app/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type StorageGet interface {
	Get(ctx context.Context, key string) (value string, ok bool)
}
type StorageList interface {
	List(ctx context.Context) (map[string]string, error)
}
type StoragePing interface {
	Ping(ctx context.Context) error
}
type StorageSet interface {
	Set(ctx context.Context, key string, value string) error
}
type StorageSetMultiple interface {
	SetMultiple(ctx context.Context, l map[string]string) error
}
type StorageDelete interface {
	Delete(ctx context.Context, l []string) error
}

type Storage interface {
	StorageGet
	StorageList
	StoragePing
	StorageSet
	StorageSetMultiple
	StorageDelete
}

func New(storage Storage, BaseURL string) chi.Router {
	r := chi.NewRouter()

	r.Mount("/debug", middleware.Profiler())

	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.TimeoutMiddleware,
			middlewares.AuthMiddleware,
			middlewares.LoggingMiddleware,
			middlewares.GzipMiddleware,
		)
		r.Get("/{id}", handlers.HandlerID(storage))
		r.Post("/*", handlers.HandlerRoot(storage, BaseURL))
		r.Post("/api/*", handlers.HandlerAPI(storage, BaseURL))
	})

	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.TimeoutMiddleware,
			middlewares.AuthMiddleware,
			middlewares.LoggingMiddleware,
		)
		r.Get("/api/user/urls", handlers.HandlerAPIUserUrls(storage, BaseURL))
		r.Delete("/api/user/urls", handlers.HandlerAPIUserUrlsDel(storage, BaseURL))
	})

	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.TimeoutMiddleware,
			middlewares.AuthMiddleware,
			middlewares.LoggingMiddleware,
			middlewares.GzipMiddleware,
		)
		r.Post("/api/shorten/batch", handlers.HandlerShortenBatch(storage, BaseURL))
	})

	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.TimeoutMiddleware,
			middlewares.LoggingMiddleware,
		)
		r.Get("/ping", handlers.HandlerPing(storage))
	})

	// empty handler for prevent error on automatic browser favicon request
	r.Get("/favicon.ico", http.HandlerFunc(handlers.HandlerEmpty))

	return r
}
