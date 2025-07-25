package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	commonconfig "github.com/rom8726/warden/internal/common/config"
)

const (
	prefix = "WARDEN"
)

type Config struct {
	Logger      commonconfig.Logger     `envconfig:"LOGGER"`
	TechServer  commonconfig.Server     `envconfig:"TECH_SERVER"`
	Postgres    commonconfig.Postgres   `envconfig:"POSTGRES"`
	ClickHouse  commonconfig.ClickHouse `envconfig:"CLICKHOUSE"`
	Mailer      commonconfig.Mailer     `envconfig:"MAILER"`
	FrontendURL string                  `default:"https://warden.your-domain" envconfig:"FRONTEND_URL"`
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
