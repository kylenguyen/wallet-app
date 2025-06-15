package server

import (
	"bitbucket.org/ntuclink/ff-order-history-go/internal/repo"
	"bitbucket.org/ntuclink/ff-order-history-go/internal/service"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	ddgin "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"

	"bitbucket.org/ntuclink/ff-order-history-go/internal/config"
	"bitbucket.org/ntuclink/ff-order-history-go/internal/handler"
)

// Server represents the HTTP server.
type Server struct {
	engine *gin.Engine
	db     *sqlx.DB
	logger *zerolog.Logger
	config config.Config
	addr   string // Add addr to the struct
}

// New creates a new HTTP server.
func New(db *sqlx.DB, logger *zerolog.Logger, cfg config.Config) *Server {
	r := gin.New()

	return &Server{
		engine: r,
		db:     db,
		logger: logger,
		config: cfg,
		addr:   fmt.Sprintf(":%d", cfg.ServicePort), // Initialize addr here
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	s.logger.Info().
		Str("addr", s.addr).
		Str("service", s.config.ServiceName).
		Str("env", s.config.Env).
		Msg("Starting HTTP server")
	s.logger.Info().Str("addr", s.addr).Msg("Starting HTTP server")

	if err := s.engine.Run(s.addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Panic().Err(err).Msg("Failed to start HTTP server")
		return err
	}

	return nil
}

// UseMiddleware adds middleware to the Gin engine.
//   - Add DataDog middleware for Gin
//   - Use Zerolog as Gin's logger
//   - Add Gin's recovery middleware
func (s *Server) UseMiddleware() {
	s.engine.Use(ddgin.Middleware(s.config.ServiceName))

	s.engine.Use(s.ginZerolog())

	s.engine.Use(gin.Recovery())
}

// ginZerolog is a middleware that logs Gin requests using Zerolog.
func (s *Server) ginZerolog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Generate or get request ID
		requestID := c.Request.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Writer.Header().Set("X-Request-Id", requestID)

		// Add request ID to logger context
		reqLogger := s.logger.With().Str("request-id", requestID).Logger()
		ctx := reqLogger.WithContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request
			for _, e := range c.Errors.Errors() {
				reqLogger.Error().Msg(e)
			}
		} else {
			// Log request details
			event := reqLogger.Info()
			if c.Writer.Status() >= http.StatusInternalServerError {
				event = reqLogger.Error()
			}

			event.
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("path", path).
				Str("query", raw).
				Str("ip", c.ClientIP()).
				Dur("latency", latency).
				Str("user-agent", c.Request.UserAgent()).
				Int("content-len", c.Writer.Size())

			if c.Writer.Status() >= http.StatusInternalServerError {
				event.Msg("Request completed with error")
			} else {
				event.Msg("Request completed")
			}
		}
	}
}

// RegisterRoutes registers the HTTP routes.
func (s *Server) RegisterRoutes() {
	// Initialize repositories, services, and handlers
	//orderSummaryRepo := repo.NewOrderSummary(s.db)
	//orderSummaryService := service.NewOrderSummary(orderSummaryRepo, datetime.NewDatetime(time.Now))
	//orderSummaryHandler := handler.NewOrderSummary(orderSummaryService)
	wRepo := repo.NewWalletImpl(s.db)
	walletService := service.NewWalletImpl(wRepo)
	walletHandler := handler.NewWalletImpl(walletService)

	s.engine.Group("/v1").
		GET("/user/:userId/wallet/:walletId/transactions", walletHandler.GetWalletTransactions)

}
