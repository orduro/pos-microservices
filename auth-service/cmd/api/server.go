package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/orduro/pos-microservices/auth-service/internal/service"
)

type server struct {
	config      *config
	userService *service.UserService
}

func (s *server) serve(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:         s.config.Addr,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server is listening on port %s\n", s.config.Addr)

	return srv.ListenAndServe()
}
