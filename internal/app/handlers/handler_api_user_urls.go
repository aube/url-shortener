package handlers

import (
	"fmt"
	"net/http"
)

func HandlerAPIUserUrls(MemoryStore Storage, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		json := string(urlsJSON(MemoryStore, baseURL))
		fmt.Fprintf(w, `%s`, json)

		w.WriteHeader(http.StatusOK)
	}
}
