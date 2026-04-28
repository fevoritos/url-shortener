package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener/config"
	infrastructurepostgres "url-shortener/internal/infrastructure/postgres"
	memoryrepo "url-shortener/internal/repository/memory"
	postgresrepo "url-shortener/internal/repository/postgres"
	transporthttp "url-shortener/internal/transport/http"
	linkhandler "url-shortener/internal/transport/http/handlers"
	"url-shortener/internal/usecase/link"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg := config.LoadConfig()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var linkRepo link.Repository

	switch cfg.StorageType {
	case config.Postgres:
		pool, err := infrastructurepostgres.Open(ctx, cfg.DatabaseDSN)
		if err != nil {
			logger.Error("open postgres", "error", err)
			os.Exit(1)
		}
		defer pool.Close()
		linkRepo = postgresrepo.New(pool)
		logger.Info("using postgres storage")

	case config.Memory:
		linkRepo = memoryrepo.New()
		logger.Info("using in-memory storage")

	default:
		logger.Error("unknown storage type", "type", cfg.StorageType)
		os.Exit(1)
	}

	linkUsecase := link.NewService(linkRepo)
	linkHandler := linkhandler.NewLinkHandler(linkUsecase)
	router := transporthttp.NewRouter(linkHandler)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown http server", "error", err)
		}
	}()

	logger.Info("http server started", "addr", cfg.HTTPAddr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("listen and serve", "error", err)
		os.Exit(1)
	}
}
