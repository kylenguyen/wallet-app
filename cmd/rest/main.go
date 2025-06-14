package main

import (
	"errors"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"bitbucket.org/ntuclink/ff-order-history-go/internal/config"
	"bitbucket.org/ntuclink/ff-order-history-go/internal/db"
	"bitbucket.org/ntuclink/ff-order-history-go/internal/server"
)

func main() {
	// Initialize Zerolog
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize DataDog tracer
	tracer.Start(
		tracer.WithEnv(cfg.Env),
		tracer.WithService(cfg.ServiceName),
	)
	defer tracer.Stop()

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
