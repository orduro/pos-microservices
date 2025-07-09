package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Addr  string     `envconfig:"SERVER_PORT" validate:"required"`
	Env   string     `envconfig:"ENVIRONMENT" validate:"required,oneof=dev development staging prod production"`
	Email *emailConf `validate:"required"`
}

type emailConf struct {
	Provider  string `envconfig:"EMAIL_PROVIDER" validate:"required"`
	FromEmail string `envconfig:"EMAIL_FROM_EMAIL" validate:"required,email"`
	FromName  string `envconfig:"EMAIL_FROM_NAME" validate:"required"`
	ReplyTo   string `envconfig:"EMAIL_REPLY_TO" validate:"required,email"`
	SMTPHost  string `envconfig:"SMTP_HOST"`
	SMTPPort  int    `envconfig:"SMTP_PORT"`
	SMTPUser  string `envconfig:"SMTP_USER"`
	SMTPPass  string `envconfig:"SMTP_PASS"`
}

func NewConfig() (*config, error) {
	var cfg config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %w", err)
	}

	// set defaults for smtp if not provided
	if cfg.Email.SMTPPort == 0 {
		cfg.Email.SMTPPort = 587
	}

	v := validator.New()
	if err := v.Struct(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
