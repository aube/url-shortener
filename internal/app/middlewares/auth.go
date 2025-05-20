package middlewares

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/aube/url-shortener/internal/app/ctxkeys"
	"github.com/aube/url-shortener/internal/logger"
)

const (
	authCookieName = "auth"    // Name of the authentication cookie
	bearerString   = "Bearer " // Prefix for bearer token in Authorization header
)

// Claims represents the JWT claims structure containing user information.
type Claims struct {
	UserID string `json:"id"` // Unique user identifier
	jwt.RegisteredClaims
}

// User represents a user entity with authentication details.
type User struct {
	Username string `json:"username"` // User's name
	Password string `json:"password"` // User's password (not currently used)
	ID       string `json:"id"`       // Unique user ID
}

// randUserID generates a random user ID between 11111 and 99999.
func randUserID() string {
	min := 11111
	max := 99999
	rndInt := rand.Intn(max-min) + min
	return strconv.Itoa(rndInt)
}

// getToken generates a new JWT token for authentication.
// The token contains a random user ID and expires in 24 hours.
func getToken() string {
	var user User
	user.ID = randUserID()

	// Create the JWT claims
	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "url-shortener-super-duper-magic-app",
		},
	}

	// Create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSecret := config.NewConfig().TokenSecret
	tokenString, err := token.SignedString(tokenSecret)

	log := logger.Get()
	if err != nil {
		log.Error("getToken", "token", token)
		log.Error("getToken", "tokenString", tokenString)
		log.Error("getToken", "claims", claims)
		return ""
	}
	return tokenString
}

// deleteAuthCookie removes the authentication cookie from the response.
func deleteAuthCookie(w http.ResponseWriter) {
	c := &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	}
	http.SetCookie(w, c)
}

// setAuthCookie sets the authentication cookie in the response.
func setAuthCookie(w http.ResponseWriter, value string) {
	log := logger.Get()
	c := &http.Cookie{
		Name:     authCookieName,
		Value:    value,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	}
	http.SetCookie(w, c)
	log.Warn("setCookie", "value", value)
}

// AuthMiddleware is a middleware that handles JWT authentication.
// It checks for a valid token in either the Authorization header or auth cookie.
// If no valid token is found, it generates a new one.
// The user ID from the token is added to the request context.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.WithContext(r.Context())

		tokenString := ""
		authHeader := r.Header.Get("Authorization")

		// Check for token in cookie if not in header
		if authHeader == "" {
			cookie, err := r.Cookie(authCookieName)
			if err == nil {
				authHeader = cookie.Value
			}
		}

		// Extract token string if header exists
		if authHeader != "" {
			tokenString = authHeader[len(bearerString):]
		}

		// Generate new token if none found
		if tokenString == "" {
			tokenString = getToken()
			w.Header().Set("Authorization", bearerString+tokenString)
		}

		// Parse and validate token
		tokenSecret := config.NewConfig().TokenSecret
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return tokenSecret, nil
		})

		tokenErrorMsg := ""
		if err != nil {
			log.Error("AuthMiddleware", "err", err)
			log.Warn("AuthMiddleware", "token", token)
			if err == jwt.ErrSignatureInvalid {
				tokenErrorMsg = "Invalid token signature"
			} else {
				tokenErrorMsg = "Invalid token"
			}
		}

		if !token.Valid {
			tokenErrorMsg = "Invalid token"
		}

		// Handle invalid token
		if tokenErrorMsg != "" {
			deleteAuthCookie(w)
			http.Error(w, tokenErrorMsg, http.StatusUnauthorized)
			return
		}

		// Set cookie and add user ID to context
		setAuthCookie(w, bearerString+tokenString)
		UserID := claims.UserID
		ctx := context.WithValue(r.Context(), ctxkeys.UserIDKey, UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
