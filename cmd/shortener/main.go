package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/aube/url-shortener/internal/app"
	"github.com/aube/url-shortener/internal/logger"
)

func main() {
	logger.Init(logger.Config{
		Level:     slog.LevelDebug,
		Output:    os.Stdout,
		AddSource: true,
		JSON:      false,
	})

	slog.SetDefault(logger.Get())

	err := app.Run()

	if err != nil {
		fmt.Println(err)
	}
}
