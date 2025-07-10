package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Addr string `envconfig:"SERVER_ADDR" validate:"required"`
	Env  string `envconfig:"ENVIRONMENT" validate:"required,oneof=dev development staging prod production"`
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
