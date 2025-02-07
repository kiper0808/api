package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/zap"

	"github.com/kiper0808/s3/internal/config"
)

type server struct {
	cfg config.HttpServer

	httpServer *http.Server

	router http.Handler

	logger *zap.Logger
}

func NewServer(logger *zap.Logger, cfg config.HttpServer, router http.Handler) *server {
	return &server{
		cfg: cfg,
		httpServer: &http.Server{
			Handler:           router,
			ReadHeaderTimeout: 0,
		},
		router: router,
		logger: logger,
	}
}

func (s *server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("create net listener failed: %w", err)
	}

	go func() {
		port := listener.Addr().(*net.TCPAddr).Port
		s.logger.Info("start http server on port %d", zap.Int("port", port))
		if err := s.httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("serve http server failed", zap.Error(err))
		}
		s.logger.Info("http server stopped")
	}()

	return nil
}

func (s *server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.CloseTimeout)
	defer cancel()

	s.httpServer.SetKeepAlivesEnabled(false)
	return s.httpServer.Shutdown(ctx)
}
