package logger

import (
	"log/slog"
	"os"
	"strings"
)

func New(level string) *slog.Logger {
	logLevel := slog.LevelInfo
	switch strings.ToUpper(level) {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	return slog.New(handler)
}
