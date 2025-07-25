package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	commonconfig "github.com/rom8726/warden/internal/common/config"
)

const (
	prefix = "WARDEN"
)

type Config struct {
	Logger     commonconfig.Logger   `envconfig:"LOGGER"`
	APIServer  commonconfig.Server   `envconfig:"API_SERVER"`
	TechServer commonconfig.Server   `envconfig:"TECH_SERVER"`
	Postgres   commonconfig.Postgres `envconfig:"POSTGRES"`
	Kafka      commonconfig.Kafka    `envconfig:"KAFKA"`
	Redis      commonconfig.Redis    `envconfig:"REDIS"`
	RateLimit  RateLimit             `envconfig:"RATE_LIMIT"`
}

type RateLimit struct {
	// RPS window in seconds
	RPSWindow time.Duration `default:"10s" envconfig:"RPS_WINDOW"`
	// Interval for refreshing RPS stats from Redis
	StatsRefreshInterval time.Duration `default:"1s" envconfig:"STATS_REFRESH_INTERVAL"`
	// Global rate limit for all projects (requests per second)
	RateLimit uint64 `default:"100" envconfig:"RATE_LIMIT"`
}

func New(filePath string) (*Config, error) {
	cfg := &Config{}

	if filePath != "" {
		if err := godotenv.Load(filePath); err != nil {
			return nil, fmt.Errorf("error loading env file: %w", err)
		}
	}

	if err := envconfig.Process(prefix, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
