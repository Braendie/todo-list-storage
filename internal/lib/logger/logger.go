package logger

import (
	"log/slog"
	"os"

	"github.com/Braendie/todo-list-storage/internal/config"
)

func MustLoad(cfg *config.Config) *slog.Logger {
	var level slog.Level
	switch cfg.Env {
	case "prod":
		level = slog.LevelInfo
	default:
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)
	return logger
}
