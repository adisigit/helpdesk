package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	JWTTTL         time.Duration
	MigrationsPath string
	CORSOrigins    []string
	SeedOnStart    bool
	SeedPassword   string
}

func Load() (*Config, error) {
	ttl, err := time.ParseDuration(getenv("JWT_TTL", "24h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_TTL: %w", err)
	}

	c := &Config{
		Port:           getenv("PORT", "8080"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTTTL:         ttl,
		MigrationsPath: getenv("MIGRATIONS_PATH", "file://db/migrations"),
		CORSOrigins:    splitCSV(getenv("CORS_ORIGINS", "http://localhost:5173")),
		SeedOnStart:    getenv("SEED_ON_START", "false") == "true",
		SeedPassword:   os.Getenv("SEED_PASSWORD"),
	}

	if c.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	if len(c.JWTSecret) < 16 {
		return nil, errors.New("JWT_SECRET is required (min 16 chars)")
	}
	if c.SeedOnStart && c.SeedPassword == "" {
		return nil, errors.New("SEED_PASSWORD is required when SEED_ON_START=true")
	}
	return c, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
