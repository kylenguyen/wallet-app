package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/kylenguyen/wallet-app/internal/config"
	"github.com/kylenguyen/wallet-app/internal/db"
	"github.com/kylenguyen/wallet-app/internal/server"
	"github.com/rs/zerolog"
)

func main() {
	// Initialize Zerolog
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Connect to the database
	db, err := db.Connect(cfg)
	if err != nil {
		panic(err)
	}

	// Initialize HTTP server
	srv := server.New(db, &logger, cfg)

	srv.UseMiddleware()

	srv.RegisterRoutes()

	if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Panic().Err(err).Msg("Failed to start HTTP server")
	}
}
