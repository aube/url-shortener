package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/aube/url-shortener/internal/app/config"
)

func TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := config.NewConfig().DefaultRequestTimeout
		ctx, cancel := context.WithTimeout(r.Context(), time.Duration(t)*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
