package main

import (
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/internal/app/handlers"
)

const portNumber = "8080"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, http.HandlerFunc(handlers.HandlerRoot))
	mux.HandleFunc(`/{id}`, http.HandlerFunc(handlers.HandlerId))

	// empty handler for automatic browser favicon request
	mux.HandleFunc(`/favicon.ico`, http.HandlerFunc(handlers.HandlerEmpty))

	fmt.Println("Server starting at:", portNumber)

	err := http.ListenAndServe(":"+portNumber, mux)

	if err != nil {
		fmt.Println("Error starting server:", err)
	} else {
		fmt.Println("Server started at:", portNumber)
	}
}
