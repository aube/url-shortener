package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func Initialize() error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	// defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugar = *logger.Sugar()

	return nil
}

func Infoln(args ...interface{}) {
	sugar.Infoln(args...)
}

func Println(args ...any) {
	fmt.Println(args...)
}
