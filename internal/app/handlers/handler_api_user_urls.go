package handlers

import (
	"fmt"
	"net/http"
)

type StorageList interface {
	List() map[string]string
}

func HandlerAPIUserUrls(MemoryStore StorageList, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		memData := MemoryStore.List()
		json := string(urlsJSON(memData, baseURL))
		fmt.Fprintf(w, `%s`, json)

		w.WriteHeader(http.StatusOK)
	}
}
