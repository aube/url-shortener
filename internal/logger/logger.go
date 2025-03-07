package logger

import (
	"fmt"
	"log/slog"
	"os"
)

// var sugar zap.SugaredLogger

func Initialize() error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	/* l, err := zap.NewDevelopment()
	   if err != nil {
	       panic(err)
	   } */

	// делаем регистратор SugaredLogger
	// sugar = *l.Sugar()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	return nil
}

func Infoln(args ...interface{}) {
	// sugar.Infoln(args...)
	slog.Info("slog.Info", args...)
}

func Println(args ...any) {
	fmt.Println(args...)
}
