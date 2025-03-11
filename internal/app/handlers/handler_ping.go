package handlers

import (
	"fmt"
	"net/http"

	"github.com/aube/url-shortener/internal/logger"
)

type StoragePing interface {
	Ping() error
}

func HandlerPing(DBStore StoragePing) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := DBStore.Ping()
		if err != nil {
			logger.Println(err)
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}

		logger.Println("Ping DB")
		fmt.Fprintf(w, `pong`)
		w.WriteHeader(http.StatusOK)
	}
}
