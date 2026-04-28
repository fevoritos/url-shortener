package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener/config"
	infrastructurepostgres "url-shortener/internal/infrastructure/postgres"
	"url-shortener/internal/lib/logger/slogpretty"
	"url-shortener/internal/middleware"
	memoryrepo "url-shortener/internal/repository/memory"
	postgresrepo "url-shortener/internal/repository/postgres"
	transporthttp "url-shortener/internal/transport/http"
	linkhandler "url-shortener/internal/transport/http/handlers"
	"url-shortener/internal/usecase/link"
	linkusecase "url-shortener/internal/usecase/link"
)

func main() {
	cfg := config.LoadConfig()
	log := slogpretty.SetupPrettySlog()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	var linkRepo linkusecase.Repository

	switch cfg.StorageType {
	case config.Postgres:
		pool, err := infrastructurepostgres.Open(ctx, cfg.DatabaseDSN)
		if err != nil {
			log.Error("open postgres", "error", err)
			os.Exit(1)
		}
		defer pool.Close()
		linkRepo = postgresrepo.New(pool)
		log.Info("using postgres storage")

	case config.Memory:
		linkRepo = memoryrepo.New()
		log.Info("using in-memory storage")

	default:
		log.Error("unknown storage type", "type", cfg.StorageType)
		os.Exit(1)
	}

	linkUsecase := link.NewService(linkRepo)
	linkHandler := linkhandler.NewLinkHandler(linkUsecase, log)
	router := transporthttp.NewRouter(linkHandler)

	stack := middleware.Chain(
		middleware.CORS,
		middleware.NewLogger(log),
	)

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           stack(router),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("shutdown http server", "error", err)
		}
	}()

	log.Info("http server started", "addr", cfg.HTTPAddr, "storage", cfg.StorageType)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("listen and serve", "error", err)
		os.Exit(1)
	}
}
