package logger

import (
	"go-vsa-template/internal/config"
	"os"

	"github.com/rs/zerolog"
)

func New(cfg *config.Config) zerolog.Logger {
	if cfg.Env == "development" {
		// На локалке (dev) делаем красиво и читаемо
		return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Logger()
	}

	return zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
}
