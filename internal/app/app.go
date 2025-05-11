package app

import (
	"log"
	"net/http"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/aube/url-shortener/internal/app/router"
	"github.com/aube/url-shortener/internal/app/store"
)

func Run() error {

	config := config.NewConfig()

	storage := store.MewStore()
	r := router.New(storage, config.BaseURL)

	address := config.ServerHost + ":" + config.ServerPort
	log.Println("Server starting", "address", address)

	err := http.ListenAndServe(address, r)

	if err != nil {
		log.Fatal("Starting server", "err", err)
	}

	return nil
}
