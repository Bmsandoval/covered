package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config is loaded from environment (see ex.env).
type Config struct {
	AppEnv            string
	AppPort           int
	AppBaseURL        string
	DatabaseDriver    string
	DatabaseDSN       string
	CacheBackend      string
	StorageBackend    string
	StorageLocalPath  string
	SessionSecret     string
	LogLevel          string
}

// Load reads configuration from the process environment.
func Load() (Config, error) {
	port, err := strconv.Atoi(envOr("APP_PORT", "8080"))
	if err != nil {
		return Config{}, fmt.Errorf("APP_PORT: %w", err)
	}

	return Config{
		AppEnv:           envOr("APP_ENV", "development"),
		AppPort:          port,
		AppBaseURL:       envOr("APP_BASE_URL", "http://localhost:8080"),
		DatabaseDriver:   envOr("DATABASE_DRIVER", "sqlite"),
		DatabaseDSN:      envOr("DATABASE_DSN", "file:./data/covered.db?cache=shared&mode=rwc"),
		CacheBackend:     envOr("CACHE_BACKEND", "memory"),
		StorageBackend:   envOr("STORAGE_BACKEND", "local"),
		StorageLocalPath: envOr("STORAGE_LOCAL_PATH", "./data/uploads"),
		SessionSecret:    envOr("SESSION_SECRET", ""),
		LogLevel:         envOr("LOG_LEVEL", "info"),
	}, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
