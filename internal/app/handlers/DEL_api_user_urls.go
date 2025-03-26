package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type StorageDelete interface {
	Delete(c context.Context, l []interface{}) error
}

func HandlerAPIUserUrlsDel(store StorageDelete, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Read the entire body content
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		var data []interface{}
		json.Unmarshal([]byte(body), &data)

		err = store.Delete(r.Context(), data)

		if err != nil {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
