package main

import (
	"context"
	"fmt"
	delivery "github.com/kiper0808/api/internal/storage/api"
	"github.com/kiper0808/api/internal/storage/config"
	"github.com/kiper0808/api/internal/storage/log"
	http3 "github.com/kiper0808/api/internal/storage/server/http"
	"github.com/kiper0808/api/internal/storage/service"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kiper0808/api/internal/storage/service/minio_storage"

	http2 "github.com/kiper0808/api/pkg/http"
	"go.uber.org/zap"
)

func main() {
	const exitFailed = 1

	fmt.Println("run karma8 storage backend") //nolint

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("failed to parse config: %s\n", err) // nolint
		os.Exit(exitFailed)
	}

	logger := log.NewLogger(cfg.LogLevel)

	if err := run(cfg, logger); err != nil {
		logger.Error("karma8 storage service: problem while trying to start / graceful shutdown server", zap.Error(err))
		os.Exit(exitFailed)
	}
}

func run(cfg *config.Config, logger *zap.Logger) error {
	// global context
	termCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// wait signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-quit
		logger.Info("signal was received", zap.Stringer("sig", sig))
		cancel()
	}()

	// init wg
	wg := &sync.WaitGroup{}

	globalHttpClient := http2.NewHTTPClient(cfg.StandardHttpClient.Timeout)

	minioStorageClient, err := minio_storage.NewClient(&cfg.MinioStorage, globalHttpClient.Client)
	if err != nil {
		logger.Fatal("create file storage client: %w", zap.Error(err))
	}

	// services & API Handlers
	services := service.NewServices(&service.Deps{
		Logger:            logger,
		HttpClient:        globalHttpClient,
		Config:            cfg,
		FileStorageClient: minioStorageClient,
	})
	handlers := delivery.NewHandler(services,
		logger,
		&cfg.MinioStorage,
		cfg.ApiVersion,
	)

	// init http servers
	httpServer := http3.NewServer(logger, cfg.Server, handlers.Init(cfg))

	if err := httpServer.Start(); err != nil {
		logger.Fatal("start http server failed", zap.Error(err))
	}

	privateHttpServer := http3.NewPrivateServer(logger, cfg.PrivateHttpServer)

	if err := privateHttpServer.Start(); err != nil {
		logger.Fatal("start private http server failed", zap.Error(err))
	}

	logger.Info("app started")

	// graceful shutdown
	<-termCtx.Done()

	if err := httpServer.Stop(); err != nil {
		logger.Error("failed to stop http server", zap.Error(err))
	}

	if err := privateHttpServer.Stop(); err != nil {
		logger.Error("failed to stop private http server", zap.Error(err))
	}

	wg.Wait()
	logger.Info("app stopped")
	return nil
}
