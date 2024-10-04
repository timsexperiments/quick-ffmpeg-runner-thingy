package log

import (
	"log/slog"
	"os"
	"sync"
)

type Logger = slog.Logger

var (
	logger *slog.Logger
	once   sync.Once
)

// InitLogger initializes the global logger with the specified log level
func InitLogger(verbose bool) {
	once.Do(func() {
		opts := &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
		if verbose {
			opts.Level = slog.LevelDebug
		}

		handler := slog.NewTextHandler(os.Stdout, opts)
		logger = slog.New(handler)
		slog.SetDefault(logger)
	})
}

// GetLogger returns the global logger instance
func L() *Logger {
	if logger == nil {
		InitLogger(false)
	}
	return logger
}
