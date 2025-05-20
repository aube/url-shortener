package router

import (
	"context"
	"net/http"

	"github.com/aube/url-shortener/internal/app/handlers"
	"github.com/aube/url-shortener/internal/app/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// StorageGet defines the interface for retrieving URL mappings.
type StorageGet interface {
	Get(ctx context.Context, key string) (value string, ok bool)
}

// StorageList defines the interface for listing URL mappings.
type StorageList interface {
	List(ctx context.Context) (map[string]string, error)
}

// StoragePing defines the interface for checking storage availability.
type StoragePing interface {
	Ping(ctx context.Context) error
}

// StorageSet defines the interface for storing URL mappings.
type StorageSet interface {
	Set(ctx context.Context, key string, value string) error
}

// StorageSetMultiple defines the interface for batch URL storage operations.
type StorageSetMultiple interface {
	SetMultiple(ctx context.Context, l map[string]string) error
}

// StorageDelete defines the interface for deleting URL mappings.
type StorageDelete interface {
	Delete(ctx context.Context, l []string) error
}

// Storage is the comprehensive interface combining all storage operations.
type Storage interface {
	StorageGet
	StorageList
	StoragePing
	StorageSet
	StorageSetMultiple
	StorageDelete
}

// New creates and configures a chi router with all application routes and middleware.
// It takes a Storage implementation and base URL as parameters and returns a configured router.
// The router is organized into several groups with different middleware combinations:
//   - Debug endpoint with profiler
//   - Main URL shortening endpoints with auth, logging, timeout and gzip middleware
//   - User URL management endpoints with auth, logging and timeout
//   - Batch operations with additional gzip support
//   - Ping endpoint without auth
func New(storage Storage, BaseURL string) chi.Router {
	r := chi.NewRouter()

	// Mount debug profiler
	r.Mount("/debug", middleware.Profiler())

	// Main URL shortening endpoints
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

	// User URL management endpoints
	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.TimeoutMiddleware,
			middlewares.AuthMiddleware,
			middlewares.LoggingMiddleware,
		)
		r.Get("/api/user/urls", handlers.HandlerAPIUserUrls(storage, BaseURL))
		r.Delete("/api/user/urls", handlers.HandlerAPIUserUrlsDel(storage, BaseURL))
	})

	// Batch operations endpoint
	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.TimeoutMiddleware,
			middlewares.AuthMiddleware,
			middlewares.LoggingMiddleware,
			middlewares.GzipMiddleware,
		)
		r.Post("/api/shorten/batch", handlers.HandlerShortenBatch(storage, BaseURL))
	})

	// Ping endpoint
	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.TimeoutMiddleware,
			middlewares.LoggingMiddleware,
		)
		r.Get("/ping", handlers.HandlerPing(storage))
	})

	// Empty handler for browser favicon requests
	r.Get("/favicon.ico", http.HandlerFunc(handlers.HandlerEmpty))

	return r
}
