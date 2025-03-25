package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

var (
	globalLogger *slog.Logger
	initOnce     sync.Once
)

// Config holds logger configuration
type Config struct {
	Level     slog.Level
	Output    io.Writer
	AddSource bool
	JSON      bool
}

// Init initializes the global logger (thread-safe)
func Init(cfg Config) {
	initOnce.Do(func() {
		// opts := &slog.HandlerOptions{
		// 	Level:     cfg.Level,
		// 	AddSource: cfg.AddSource,
		// }

		opts := &slog.HandlerOptions{
			Level:     cfg.Level,
			AddSource: cfg.AddSource,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Customize attribute output
				if a.Key == slog.TimeKey {
					return slog.Attr{
						Key:   "ts",
						Value: slog.StringValue(a.Value.Time().Format(time.RFC3339)),
					}
				}
				return a
			},
		}

		var handler slog.Handler
		if cfg.JSON {
			handler = slog.NewJSONHandler(cfg.Output, opts)
		} else {
			handler = slog.NewTextHandler(cfg.Output, opts)
		}

		globalLogger = slog.New(handler)
	})
}

/* func caller() slog.Attr {
	_, file, line, _ := runtime.Caller(2)
	return slog.String("caller", file+":"+strconv.Itoa(line))
}

log.Info("starting processing",
	"input", input,
	logger.caller(),
)
*/

// Get returns the global logger instance
func Get() *slog.Logger {
	if globalLogger == nil {
		Init(Config{
			Level:     slog.LevelInfo,
			Output:    os.Stdout,
			AddSource: false,
			JSON:      true,
		})
	}
	return globalLogger
}

// WithContext creates a new logger with context values
func WithContext(ctx context.Context) *slog.Logger {
	logger := Get()

	// Add request ID if present in context
	if reqID, ok := ctx.Value("request_id").(string); ok {
		logger = logger.With("request_id", reqID)
	}

	return logger
}
