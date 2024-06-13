package main

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"online-store-api/internal/app/config"
	"os"
	"os/signal"
	"time"
)

func main() {
	e := echo.New()
	cfg := config.NewViperConfig()
	log := config.NewLogger(cfg)
	validator := config.NewValidator()
	db, err := config.NewPostgreDatabase(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start, error connect to DB Postgre")
		return
	}

	config.BootstrapApp(&config.BootstrapAppConfig{
		DB:        db,
		Validator: validator,
		Config:    cfg,
		Echo:      e,
	})
	startServer(e)
}

func startServer(e *echo.Echo) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
