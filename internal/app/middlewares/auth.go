package middlewares

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/aube/url-shortener/internal/app/ctxkeys"
	"github.com/aube/url-shortener/internal/logger"
)

const authCookieName = "auth"

func authenticateUser() bool {
	return true
}

func deleteAuthCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0), // Cookie expires in 24 hours
		Path:     "/",             // Cookie is accessible across the entire site
		HttpOnly: true,            // Cookie is not accessible via JavaScript
		Secure:   false,           // Set to true if using HTTPS
	}

	http.SetCookie(w, c)
}

func setAuthCookie(w http.ResponseWriter, value string) {
	log := logger.Get()
	c := &http.Cookie{
		Name:     authCookieName,
		Value:    value,
		Expires:  time.Now().Add(24 * time.Hour), // Cookie expires in 24 hours
		Path:     "/",                            // Cookie is accessible across the entire site
		HttpOnly: true,                           // Cookie is not accessible via JavaScript
		Secure:   false,                          // Set to true if using HTTPS
	}

	http.SetCookie(w, c)
	log.Warn("setCookie", "value", value)
}

func randUserID() string {
	min := 11111
	max := 99999
	rndInt := rand.Intn(max-min) + min
	return strconv.Itoa(rndInt)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Authorization")
		if userID == "" {
			cookie, err := r.Cookie(authCookieName)
			if err == http.ErrNoCookie {
				userID = randUserID()
			} else {
				userID = cookie.Value
			}
		}
		w.Header().Set("Authorization", userID)
		setAuthCookie(w, userID)

		ctx := context.WithValue(r.Context(), ctxkeys.UserIDKey, userID)
		newReq := r.WithContext(ctx)

		if !authenticateUser() {
			deleteAuthCookie(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, newReq)
	})
}
