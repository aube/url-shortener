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

// Config holds the configuration for initializing the global logger.
type Config struct {
	Level     slog.Level // Level of logging verbosity
	Output    io.Writer  // Output destination of logs
	AddSource bool       // Include source file and line number in log statements
	JSON      bool       // Whether to output logs in JSON format
}

// Init initializes the global logger. This is thread-safe and idempotent, meaning it can be called multiple times without changing anything.
func Init(cfg Config) {
	initOnce.Do(func() {
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

// Get returns a reference to the global slog.Logger instance.
// If the global logger has not been initialized before, it will be automatically initialized
// using default settings (info level logging, output directed to os.Stdout,
// source files are included, logs in JSON format).
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

// WithContext creates a new logger with context values from the provided context.
// If there is a "request_id" key in the context, it will be added to the logger's attributes.
func WithContext(ctx context.Context) *slog.Logger {
	logger := Get()

	// Add request ID if present in context
	if reqID, ok := ctx.Value("request_id").(string); ok {
		logger = logger.With("request_id", reqID)
	}

	return logger
}
