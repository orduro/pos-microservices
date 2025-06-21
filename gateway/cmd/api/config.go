package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	addr string
	env  string
}

func NewConfig() *config {
	cfg := config{
		addr: os.Getenv("SERVER_ADDR"),
		env:  os.Getenv("ENVIRONMENT"),
	}

	return &cfg
}
