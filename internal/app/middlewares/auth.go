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

const authCookieName = "auth"
const bearerString = "Bearer "

type Claims struct {
	UserID string `json:"id"`
	jwt.RegisteredClaims
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       string `json:"id"`
}

func randUserID() string {
	min := 11111
	max := 99999
	rndInt := rand.Intn(max-min) + min
	return strconv.Itoa(rndInt)
}

func getToken() string {
	var user User
	// user.Username = "admin"
	// user.Password = "password"
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

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSecret := config.NewConfig().TokenSecret

	// Sign the token with our secret
	tokenString, err := token.SignedString(tokenSecret)

	log := logger.Get()

	if err != nil {
		log.Error("getToken", "token", token)
		log.Error("getToken", "tokenString", tokenString)
		log.Error("getToken", "claims", claims)
		// http.Error(w, "Error generating token", http.StatusInternalServerError)
		return ""
	}
	return tokenString
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := ""
		authHeader := r.Header.Get("Authorization")

		log := logger.WithContext(r.Context())

		if authHeader != "" {
			// The token should be in the format "Bearer <token>"
			tokenString = authHeader[len(bearerString):]
		}

		if tokenString == "" {
			tokenString = getToken()
			w.Header().Set("Authorization", bearerString+tokenString)
		}

		tokenSecret := config.NewConfig().TokenSecret

		// Parse the token
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

		if tokenErrorMsg != "" {
			deleteAuthCookie(w)
			http.Error(w, tokenErrorMsg, http.StatusUnauthorized)
			return
		}

		setAuthCookie(w, bearerString+tokenString)
		UserID := claims.UserID

		ctx := context.WithValue(r.Context(), ctxkeys.UserIDKey, UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
