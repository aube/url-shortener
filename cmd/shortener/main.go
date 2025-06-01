package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/aube/url-shortener/internal/app"
	"github.com/aube/url-shortener/internal/logger"
)

var (
	buildVersion string
	buildTime    string
	buildCommit  string
)

func main() {
	fmt.Printf("Build version: %s\n", naOnEmpty(buildVersion))
	fmt.Printf("Build date: %s\n", naOnEmpty(buildTime))
	fmt.Printf("Build commit: %s\n\n", naOnEmpty(buildCommit))

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

func naOnEmpty(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}
