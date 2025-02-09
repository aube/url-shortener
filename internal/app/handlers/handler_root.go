package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aube/url-shortener/internal/app/hashes"
)

func HandlerRoot(w http.ResponseWriter, r *http.Request, linkAddress string) {
	switch r.Method {
	case "POST":
		// Read the entire body content
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			break
		}
		defer r.Body.Close()

		if len(body) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			break
		}

		w.WriteHeader(http.StatusCreated)
		hash := hashes.SetURLHash(body)

		fmt.Fprintf(w, linkAddress+"/"+hash)
		fmt.Println("URL:", linkAddress+"/"+hash)
	default:
		fmt.Println("Not served method:", r.Method)
	}
}
