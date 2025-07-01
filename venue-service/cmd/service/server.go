package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/orduro/pos-microservices/venue/internal/store"
)

type server struct {
	config *config
	store  store.Store
}

func (s *server) serve(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:         s.config.Addr,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("service is listening on port %s\n", s.config.Addr)

	return srv.ListenAndServe()
}
