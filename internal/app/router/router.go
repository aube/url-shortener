package router

import (
	"net/http"

	"github.com/aube/url-shortener/internal/app/handlers"
	"github.com/aube/url-shortener/internal/app/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New creates and configures a chi router with all application routes and middleware.
// It takes a Storage implementation and base URL as parameters and returns a configured router.
// The router is organized into several groups with different middleware combinations:
//   - Debug endpoint with profiler
//   - Main URL shortening endpoints with auth, logging, timeout and gzip middleware
//   - User URL management endpoints with auth, logging and timeout
//   - Batch operations with additional gzip support
//   - Ping endpoint without auth
func New(BaseURL string) chi.Router {
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
		r.Get("/{id}", handlers.HandlerID())
		r.Post("/*", handlers.HandlerRoot(BaseURL))
		r.Post("/api/*", handlers.HandlerAPI(BaseURL))
	})

	// User URL management endpoints
	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.TimeoutMiddleware,
			middlewares.AuthMiddleware,
			middlewares.LoggingMiddleware,
		)
		r.Get("/api/user/urls", handlers.HandlerAPIUserUrls(BaseURL))
		r.Delete("/api/user/urls", handlers.HandlerAPIUserUrlsDel(BaseURL))
	})

	// Batch operations endpoint
	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.TimeoutMiddleware,
			middlewares.AuthMiddleware,
			middlewares.LoggingMiddleware,
			middlewares.GzipMiddleware,
		)
		r.Post("/api/shorten/batch", handlers.HandlerShortenBatch(BaseURL))
	})

	// Ping endpoint
	r.Group(func(r chi.Router) {
		r.Use(
			middlewares.TimeoutMiddleware,
			middlewares.LoggingMiddleware,
		)
		r.Get("/ping", handlers.HandlerPing())
		r.Get("/internal/stats", handlers.HandlerInternalStats())
	})

	// Empty handler for browser favicon requests
	r.Get("/favicon.ico", http.HandlerFunc(handlers.HandlerEmpty))

	return r
}
