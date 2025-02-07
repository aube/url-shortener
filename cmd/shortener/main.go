package main

import (
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/internal/app/handlers"
	"github.com/go-chi/chi/v5"
)

const portNumber = "8080"

func main() {
	r := chi.NewRouter()

	r.Post("/", http.HandlerFunc(handlers.HandlerRoot))
	r.Get("/{id}", http.HandlerFunc(handlers.HandlerId))

	// empty handler for automatic browser favicon request
	r.Get("/favicon.ico", http.HandlerFunc(handlers.HandlerEmpty))

	fmt.Println("Server starting at:", portNumber)

	err := http.ListenAndServe(":"+portNumber, r)

	if err != nil {
		fmt.Println("Error starting server:", err)
	} else {
		fmt.Println("Server started at:", portNumber)
	}
}
