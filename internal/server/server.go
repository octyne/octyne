package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/octyne/octyne/internal/gateway"
	"github.com/octyne/octyne/internal/registry"
)

const shutdownTimeout = 30 * time.Second

type Server struct {
	mux           *http.ServeMux
	gateway       *gateway.Service
	modelRegistry *registry.Registry
	httpServer    *http.Server
	logger        *slog.Logger
}

func New(addr string, logger *slog.Logger, gateway *gateway.Service, modelRegistry *registry.Registry) *Server {

	s := &Server{
		mux:           http.NewServeMux(),
		gateway:       gateway,
		modelRegistry: modelRegistry,
		logger:        logger,
	}

	s.routes()

	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           withRequestID(s.withRequestLogging(s.mux)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      0, // Long-lived SSE streams must not have a global write deadline.
		IdleTimeout:       120 * time.Second,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	return s
}

func (s *Server) Start() error {
	s.logger.Info("HTTP server starting", "address", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Run(ctx context.Context) error {
	serverErr := make(chan error, 1)

	go func() {
		serverErr <- s.Start()
	}()

	select {
	case err := <-serverErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve HTTP server: %w", err)

	case <-ctx.Done():
		s.logger.Info("HTTP server shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		_ = s.httpServer.Close()
		return fmt.Errorf("shut down HTTP server: %w", err)
	}

	if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP server: %w", err)
	}

	return nil
}
