package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"

	"github.com/0xrinful/reddit-clone/internal/config"
)

func New(cfg config.LoggingConfig) *slog.Logger {
	if cfg.IsProduction {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{TimeFormat: "2006/01/02 15:04:05"}))
}
