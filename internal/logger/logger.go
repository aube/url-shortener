package logger

import (
	"log/slog"
	"os"
)

func Initialize() error {
	logger := slog.New(
		slog.NewJSONHandler(os.Stderr,
			&slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	return nil
}

func Warnln(args ...interface{}) {
	slog.Warn("*", args...)
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
