package logger

import (
	"log/slog"
	"os"
)

func Initialize() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	return nil
}

func Infoln(args ...interface{}) {
	slog.Info("slog.Info", args...)
}

func Errorln(args ...interface{}) {
	slog.Error("slog.Error", args...)
}

func Println(args ...any) {
	slog.Info("slog.Println", args...)
}
