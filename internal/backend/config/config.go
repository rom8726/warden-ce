//nolint:lll // it's ok
package config

import (
	"fmt"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2" //nolint:revive // it's ok
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	commonconfig "github.com/rom8726/warden/internal/common/config"
)

const (
	prefix = "WARDEN"
)

type Config struct {
	Logger           commonconfig.Logger     `envconfig:"LOGGER"`
	APIServer        commonconfig.Server     `envconfig:"API_SERVER"`
	TechServer       commonconfig.Server     `envconfig:"TECH_SERVER"`
	Postgres         commonconfig.Postgres   `envconfig:"POSTGRES"`
	ClickHouse       commonconfig.ClickHouse `envconfig:"CLICKHOUSE"`
	Mailer           commonconfig.Mailer     `envconfig:"MAILER"`
	SecretKey        string                  `envconfig:"SECRET_KEY"                         required:"true"`
	JWTSecretKey     string                  `envconfig:"JWT_SECRET_KEY"                     required:"true"`
	AccessTokenTTL   time.Duration           `default:"3h"                                   envconfig:"ACCESS_TOKEN_TTL"`
	RefreshTokenTTL  time.Duration           `default:"168h"                                 envconfig:"REFRESH_TOKEN_TTL"`
	ResetPasswordTTL time.Duration           `default:"8h"                                   envconfig:"RESET_PASSWORD_TTL"`
	LogoURL          string                  `default:"https://warden-project.tech/logo.png" envconfig:"LOGO_URL"`
	FrontendURL      string                  `default:"https://warden.your-domain"           envconfig:"FRONTEND_URL"`
	AdminEmail       string                  `envconfig:"ADMIN_EMAIL"`
	AdminTmpPassword string                  `envconfig:"ADMIN_TMP_PASSWORD"`
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
