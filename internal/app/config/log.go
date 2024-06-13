package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

func NewLogger(cfg *viper.Viper) (zrLog zerolog.Logger) {
	log.Logger = log.With().Caller().Logger()
	zerolog.TimeFieldFormat = time.RFC3339

	switch cfg.GetString("LOG_LEVEL") {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zrLog = log.Logger

	return
}
