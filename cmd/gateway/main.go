package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	delivery "github.com/kiper0808/api/internal/gateway/api"
	"github.com/kiper0808/api/internal/gateway/config"
	"github.com/kiper0808/api/internal/gateway/db"
	"github.com/kiper0808/api/internal/gateway/log"
	"github.com/kiper0808/api/internal/gateway/repository"
	http3 "github.com/kiper0808/api/internal/gateway/server/http"
	"github.com/kiper0808/api/internal/gateway/service"
	"github.com/kiper0808/api/internal/gateway/service/file_storage"

	http2 "github.com/kiper0808/api/pkg/http"
	"go.uber.org/zap"
)

func main() {
	const exitFailed = 1

	fmt.Println("run karma8 gateway backend") //nolint

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("failed to parse config: %s\n", err) // nolint
		os.Exit(exitFailed)
	}

	logger := log.NewLogger(cfg.LogLevel)

	if err := run(cfg, logger); err != nil {
		logger.Error("karma8 gateway service: problem while trying to start / graceful shutdown server", zap.Error(err))
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

	// init db
	dbMysql, err := db.New(cfg.Database)
	if err != nil {
		logger.Fatal("mysql connect problem: %w", zap.Error(err))
	}
	defer dbMysql.Close()

	globalHttpClient := http2.NewHTTPClient(cfg.StandardHttpClient.Timeout)

	fileStorageClient, err := file_storage.NewClient(globalHttpClient.Client, cfg.FileStorage)
	if err != nil {
		logger.Fatal("create file storage client: %w", zap.Error(err))
	}

	// services, repos & API Handlers
	repos := repository.NewRepositories(dbMysql, logger)

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

	go memStats()

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

var maxMemAlloc uint64 = 20 * 1024 * 1024 // 50 MB - initial threshold

func memStats() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		// Dump only if current allocation exceeds previous maximum
		if m.Alloc > maxMemAlloc {
			filename := fmt.Sprintf("/tmp/heap_dump_%d.prof", m.Alloc)
			f, err := os.Create(filename)
			if err != nil {
				slog.Error("could not create file: ", "err", err)
				continue
			}
			err = pprof.WriteHeapProfile(f)
			if err != nil {
				slog.Error("could not write heap profile: ", "err", err)
				continue
			}
			slog.Info("heap profile written - new memory peak reached",
				"filename", filename,
				"Alloc", m.Alloc,
				"TotalAlloc", m.TotalAlloc,
				"NumGoroutine", runtime.NumGoroutine(),
				"HeapSys", m.HeapSys,
				"HeapAlloc", m.HeapAlloc)

			// Update maximum to current allocation
			maxMemAlloc = m.Alloc

			if err = f.Close(); err != nil {
				slog.Error("could not close file: ", "err", err)
			}
		}
		time.Sleep(700 * time.Millisecond)
	}
}
