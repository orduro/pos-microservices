package main

import "log"

func main() {
	// get config from .env
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := &server{
		config: cfg,
	}

	// routes using chi router (routes.go)
	mux := app.mount()

	// start server
	log.Fatal(app.serve(mux))
}
