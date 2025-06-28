package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type application struct {
	config *config
}

func (app *application) serve(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("venue service is listening on port %s\n", app.config.Addr)

	return srv.ListenAndServe()
}
