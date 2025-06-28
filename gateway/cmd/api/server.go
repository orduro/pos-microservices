package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type server struct {
	config *config
}

func (s *server) serve(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:         s.config.Addr,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("gateway server is listening on port %s\n", s.config.Addr)

	return srv.ListenAndServe()
}
