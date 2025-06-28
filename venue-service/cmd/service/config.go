package main

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	Addr string `validate:"required"`
	Env  string `validate:"required,oneof=dev development staging prod production"`
}

func NewConfig() (*config, error) {
	cfg := config{
		Addr: os.Getenv("SERVER_PORT"),
		Env:  os.Getenv("ENVIRONMENT"),
	}

	v := validator.New()
	if err := v.Struct(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
