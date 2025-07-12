package main

import (
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

	srv := &server{
		config: cfg,
	}

	// routes using chi router (routes.go)
	mux := srv.mount()

	// start server
	log.Fatal(srv.serve(mux))
}
