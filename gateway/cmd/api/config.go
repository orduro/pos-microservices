package main

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	Addr string `validate:"required"`
	Env  string `validate:"required"`
}

func NewConfig() (*config, error) {
	cfg := config{
		Addr: os.Getenv("SERVER_ADDR"),
		Env:  os.Getenv("ENVIRONMENT"),
	}

	v := validator.New()
	if err := v.Struct(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	if err := validateEnvironment(cfg.Env); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateEnvironment(env string) error {
	validEnvs := []string{"dev", "development", "staging", "prod", "production"}

	for _, validEnv := range validEnvs {
		if env == validEnv {
			return nil
		}
	}

	return fmt.Errorf("invalid ENVIRONMENT: %s. Must be one of: %v", env, validEnvs)
}
