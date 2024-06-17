package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/http"
	"online-store-api/internal/app/config"
	"os"
	"os/signal"
	"time"
)

func main() {
	e := echo.New()
	e.Debug = true

	cfg := config.NewViperConfig()
	log := config.NewLogger(cfg)
	validator := config.NewValidator()
	db, err := config.NewPostgreDatabase(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start, error connect to DB Postgre")
		return
	}
	redisClient := config.NewRedisClient(cfg)

	config.BootstrapApp(&config.BootstrapAppConfig{
		DB:        db,
		Validator: validator,
		Config:    cfg,
		Echo:      e,
		Cache:     redisClient,
	})
	config.SetCustomErrorHandler(e)
	startServer(e, cfg)
}

func startServer(e *echo.Echo, cfg *viper.Viper) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		appPort := fmt.Sprintf(":%s", cfg.GetString("APP_PORT"))
		log.Info().Msgf("Starting %s server on port %s", cfg.GetString("APP_NAME"), appPort)
		if err := e.Start(appPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Msg("shutting down the server")
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Send()
	}
}
