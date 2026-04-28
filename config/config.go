package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

type StorageType string

const (
	Postgres StorageType = "postgres"
	Memory   StorageType = "memory"
)

type Config struct {
	HTTPAddr    string
	DatabaseDSN string
	StorageType StorageType
}

func LoadConfig() Config {
	_ = godotenv.Load()

	var (
		fAddr    string
		fDSN     string
		fStorage string
	)

	flag.StringVar(&fAddr, "addr", ":8080", "HTTP server address")
	flag.StringVar(&fDSN, "db-dsn", "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable", "PostgreSQL DSN")
	flag.StringVar(&fStorage, "storage", "postgres", "Storage type: postgres or memory")

	flag.Parse()

	cfg := Config{
		HTTPAddr:    envOrDefault("HTTP_ADDR", fAddr),
		DatabaseDSN: envOrDefault("DATABASE_DSN", fDSN),
		StorageType: StorageType(envOrDefault("STORAGE_TYPE", fStorage)),
	}

	return cfg
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
