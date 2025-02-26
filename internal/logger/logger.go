package logger

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func Initialize() error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugar = *logger.Sugar()

	return nil
}

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func (r *loggingResponseWriter) String() {
	var buf bytes.Buffer

	buf.WriteString("Response:")

	buf.WriteString("Headers:")
	for k, v := range r.ResponseWriter.Header() {
		buf.WriteString(fmt.Sprintf("%s: %v", k, v))
	}
}

func LoggingMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		sugar.Infoln(
			responseData.status, // получаем перехваченный код статуса ответа
			r.Method,
			r.RequestURI,
			duration,
			responseData.size, // получаем перехваченный размер ответа
			"ce", r.Header.Get("Content-Encoding"),
			"ct", r.Header.Get("Content-Type"),
			"-", w.Header(),
		)
	}
}
