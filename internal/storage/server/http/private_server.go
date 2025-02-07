package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.uber.org/zap"

	"github.com/kiper0808/s3/internal/config"
	"github.com/kiper0808/s3/internal/healthcheck"
)

type privateServer struct {
	cfg config.HttpServer

	httpServer *http.Server

	logger *zap.Logger
}

func NewPrivateServer(logger *zap.Logger, cfg config.HttpServer) *privateServer {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/health/live", healthcheck.NewHealthCheck().LiveHandler())
	router.GET("/health/ready", healthcheck.NewHealthCheck().ReadyHandler())
	router.GET("/metrics", func(c *gin.Context) {
		handler := promhttp.Handler()
		handler.ServeHTTP(c.Writer, c.Request)
	})

	return &privateServer{
		cfg: cfg,
		httpServer: &http.Server{
			Handler:           router,
			ReadHeaderTimeout: 0,
		},
		logger: logger,
	}
}

func (p *privateServer) Start() error {
	privateListener, err := net.Listen("tcp", fmt.Sprintf(":%d", p.cfg.DebugPort))
	if err != nil {
		return fmt.Errorf("create private net listener failed: %w", err)
	}

	go func() {
		port := privateListener.Addr().(*net.TCPAddr).Port
		p.logger.Info("start http private server on port %d", zap.Int("port", port))
		if err := p.httpServer.Serve(privateListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			p.logger.Fatal("serve private http private server failed", zap.Error(err))
		}
		p.logger.Info("http private server stopped")
	}()

	return nil
}

func (p *privateServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), p.cfg.CloseTimeout)
	defer cancel()

	p.httpServer.SetKeepAlivesEnabled(false)
	return p.httpServer.Shutdown(ctx)
}
