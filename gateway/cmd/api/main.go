package main

import (
	_ "github.com/orduro/pos-microservices/gateway/cmd/api/docs"
	"log"
)

// @title Gateway API
// @version 1.0
// @description This is the gateway that client's will use to access backend resources and services

// @contact.name Austin Sofaer (Backend Engineer)
// @contact.url www.linkedin.com/in/austinsofaer
// @contact.email austinsofaer@gmail.com

// @BasePath http://localhost:8080
func main() {
	// get config from .env
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		config: cfg,
	}

	// routes using chi router (routes.go)
	mux := app.mount()

	// start server
	log.Fatal(app.serve(mux))
}
