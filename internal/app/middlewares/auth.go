package middlewares

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type userID string

const userIDKey = userID("userID")

func authenticateUser(userName, password string) bool {
	// if password == strings.ToUpper(userName) {
	// 	return true
	// }
	return true
}

func deleteCookie(w http.ResponseWriter, name string) {
	c := &http.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Unix(0, 0), // Cookie expires in 24 hours
		Path:     "/",             // Cookie is accessible across the entire site
		HttpOnly: true,            // Cookie is not accessible via JavaScript
		Secure:   false,           // Set to true if using HTTPS
	}

	http.SetCookie(w, c)
}

func setCookie(w http.ResponseWriter, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  time.Now().Add(24 * time.Hour), // Cookie expires in 24 hours
		Path:     "/",                            // Cookie is accessible across the entire site
		HttpOnly: true,                           // Cookie is not accessible via JavaScript
		Secure:   false,                          // Set to true if using HTTPS
	}

	http.SetCookie(w, c)
}

func randUserID() string {
	min := 11111
	max := 99999
	rndInt := rand.Intn(max-min) + min
	return strconv.Itoa(rndInt)
}

func AuthMiddleware(next http.Handler) http.Handler {
	// deleteCookie(w, "auth")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")

		var userID string

		if err == http.ErrNoCookie {
			userID = randUserID()
			setCookie(w, "auth", userID)
		} else {
			userID = cookie.Value
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		r = r.WithContext(ctx)

		username := "123"
		password := "123"
		if !authenticateUser(username, password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// package main

// import (
// 	"net/http"
// 	"strings"
// // )

// func main() {
// 	authorizeAdmin := authMiddleware.NewAuthorizationMiddleware(func(r *http.Request, value string) bool {
// 		return value == "admin"
// 	})
// 	http.Handle("/", helloWorld)
// 	http.Handle("/admin", authenticate(authorizeAdmin(helloAdmin)))
// 	err := http.ListenAndServe("localhost:9090", nil)
// 	if err != nil {
// 		panic(err)
// 	}
// }
