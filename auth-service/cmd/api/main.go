package main

import (
	"log"
)

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
