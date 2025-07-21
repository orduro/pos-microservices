package main

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Addr     string        `envconfig:"SERVER_PORT" validate:"required"`
	Env      string        `envconfig:"ENVIRONMENT" validate:"required,oneof=dev development staging prod production"`
	Postgres *postgresConf `validate:"required"`
	Services *servicesConf `validate:"required"`
	JWT      *jwtConf      `validate:"required"`
}

type postgresConf struct {
	DSN          string        `envconfig:"POSTGRES_DSN" validate:"required"`
	MaxOpenConns int           `envconfig:"POSTGRES_MAX_OPEN_CONNS" validate:"required,min=1"`
	MaxIdleConns int           `envconfig:"POSTGRES_MAX_IDLE_CONNS" validate:"required,min=0"`
	MaxIdleTime  time.Duration `envconfig:"POSTGRES_MAX_IDLE_TIME" validate:"required"`
}

type servicesConf struct {
	BaseURL                 string `envconfig:"BASE_URL" validate:"required,url"`
	FrontendAdminURL        string `envconfig:"FRONTEND_ADMIN_URL" validate:"required,url"`
	CommunicationServiceURL string `envconfig:"COMMUNICATION_SERVICE_URL" validate:"required"`
}

// token defaults:
// - 3 min access token expiration
// - 7 days refresh token expiration
type jwtConf struct {
	Secret                string        `envconfig:"JWT_SECRET" validate:"required,min=32"`
	ExpirationTime        time.Duration `envconfig:"JWT_EXPIRATION_TIME" default:"3m"`
	RefreshExpirationTime time.Duration `envconfig:"JWT_REFRESH_EXPIRATION_TIME" default:"168h"`
}

func NewConfig() (*config, error) {
	var cfg config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}

	v := validator.New()
	if err := v.Struct(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
