package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	hashes "github.com/aube/url-shortener/internal/app/hashes"
)

const portNumber = "8080"

func HandlerRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		fmt.Println(r.Method)

		// Read the entire body content
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		hash := hashes.SetURLHash(body)

		fmt.Fprintf(w, "http://localhost:8080/"+hash)
		fmt.Println("URL:", "http://localhost:8080/"+hash)
	default:
		fmt.Println("Not served method:", r.Method)
	}
}
