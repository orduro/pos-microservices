package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type config struct {
	Addr     string        `envconfig:"SERVER_PORT" validate:"required"`
	Env      string        `envconfig:"ENVIRONMENT" validate:"required,oneof=dev development staging prod production"`
	Postgres *postgresConf `validate:"required"`
}

type postgresConf struct {
	DSN          string        `envconfig:"POSTGRES_DSN" validate:"required"`
	MaxOpenConns int           `envconfig:"POSTGRES_MAX_OPEN_CONNS" validate:"required,min=1"`
	MaxIdleConns int           `envconfig:"POSTGRES_MAX_IDLE_CONNS" validate:"required,min=0"`
	MaxIdleTime  time.Duration `envconfig:"POSTGRES_MAX_IDLE_TIME" validate:"required"`
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
