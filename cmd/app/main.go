package main

import (
	"context"
	"fmt"
	delivery "github.com/kiper0808/api/internal/gateway/api"
	"github.com/kiper0808/api/internal/gateway/cache"
	"github.com/kiper0808/api/internal/gateway/config"
	"github.com/kiper0808/api/internal/gateway/db"
	"github.com/kiper0808/api/internal/gateway/log"
	"github.com/kiper0808/api/internal/gateway/metrics"
	"github.com/kiper0808/api/internal/gateway/repository"
	http3 "github.com/kiper0808/api/internal/gateway/server/http"
	"github.com/kiper0808/api/internal/gateway/service"
	"github.com/kiper0808/api/internal/gateway/service/file_storage"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/extra/redisprometheus/v9"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	http2 "github.com/kiper0808/api/pkg/http"

	"go.uber.org/zap"
)

func main() {
	const exitFailed = 1

	fmt.Println("run dostavkee backend") //nolint

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("failed to parse config: %s\n", err) // nolint
		os.Exit(exitFailed)
	}

	logger := log.NewLogger(cfg.LogLevel)

	if err := run(cfg, logger); err != nil {
		logger.Error("dostavkee service: problem while trying to start / graceful shutdown server", zap.Error(err))
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

	// init redis cache
	redis, err := cache.NewRedis(cfg.Cache)
	if err != nil {
		return fmt.Errorf("redis init problem: %w", err)
	}

	if err := redisotel.InstrumentTracing(redis); err != nil {
		return fmt.Errorf("redis tracing init problem: %w", err)
	}

	// init db
	dbMysql, err := db.New(cfg.Database)
	if err != nil {
		logger.Fatal("mysql connect problem: %w", zap.Error(err))
	}
	defer dbMysql.Close()

	dbCollector := collectors.NewDBStatsCollector(dbMysql.DB, cfg.Database.DBName)
	redisCollector := redisprometheus.NewCollector("dostavkee", "redis", redis)

	if err := registerCollectors(dbCollector, redisCollector); err != nil {
		logger.Fatal("register collectors failed: %w", zap.Error(err))
	}

	if err := metrics.RegisterMetrics(); err != nil {
		logger.Fatal("register metrics problem: %w", zap.Error(err))
	}

	globalHttpClient := http2.NewHTTPClient(cfg.StandardHttpClient.Timeout)

	fileStorageClient, err := fileStorage.NewClient(globalHttpClient.Client, cfg.FileStorage)
	if err != nil {
		logger.Fatal("create file storage client: %w", zap.Error(err))
	}

	// services, workers, repos & API Handlers
	repos := repository.NewRepositories(dbMysql, redis, logger)

	services := service.NewServices(&service.Deps{
		Logger:            logger,
		HttpClient:        globalHttpClient,
		Repos:             repos,
		Config:            cfg,
		FileStorageClient: fileStorageClient,
	})
	handlers := delivery.NewHandler(services,
		logger,
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

func registerCollectors(promCollectors ...prometheus.Collector) error {
	for i := range promCollectors {
		if err := prometheus.Register(promCollectors[i]); err != nil {
			return fmt.Errorf("prometheus collector register fail: %w", err)
		}
	}
	return nil
}
