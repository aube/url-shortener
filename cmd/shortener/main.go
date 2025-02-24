package main

import (
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/internal/app/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	serverAddress, baseURL := NewConfig()

	r := chi.NewRouter()

	r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerRoot(w, r, baseURL)
	})
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlerID(w, r)
	})

	// empty handler for prevent error on automatic browser favicon request
	r.Get("/favicon.ico", http.HandlerFunc(handlers.HandlerEmpty))

	err := http.ListenAndServe(serverAddress, r)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
