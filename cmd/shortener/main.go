package main

import (
	"fmt"

	"github.com/aube/url-shortener/internal/app"
)

func main() {
	err := app.Run()

	if err != nil {
		fmt.Println(err)
	}
}
