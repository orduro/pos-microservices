package main

import (
	"log"
)

func main() {
	app := &application{
		config: NewConfig(),
	}

	mux := app.mount()

	// create server and run
	log.Fatal(app.serve(mux))
}
