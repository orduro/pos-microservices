package main

import (
	"log"
	"net/http"
)

func main() {
	// multiplexer configs in routes.go
	r := getMux()

	// create server and run
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Printf("Gateway server started at localhost%v", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Gateway server encountered an error: %v", err)
	}
}
