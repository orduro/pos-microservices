package main

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/orduro/common/auth"
	"github.com/orduro/pos-microservices/venue/internal/handler/health"
	"github.com/orduro/pos-microservices/venue/internal/handler/venue"
)

func (s *server) mount() *chi.Mux {
	// initialise handlers
	healthhandler := health.New(s.config.Env)
	venuehandler := venue.New(s.store)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", healthhandler.Check)

	// prefix /api in front of all routes, all routes go in here
	r.Route("/api", func(r chi.Router) {
		// make sure secret is correct for /api access
		r.Use(auth.JWTAuth(s.config.JWT.Secret))

		// routes for /venue
		r.Route("/venue", func(r chi.Router) {
			r.Post("/", venuehandler.CreateVenue)
			r.Get("/{id}", venuehandler.GetVenueByID)
			r.Patch("/{id}", venuehandler.UpdateVenue)

			// routes for venue archives `/venue/archive`
			r.Route("/archive", func(r chi.Router) {
				r.Put("/{id}", venuehandler.ArchiveVenue)
			})

			// routes for venue deletes `/venue/delete`
			r.Route("/delete", func(r chi.Router) {
				r.Delete("/{id}", venuehandler.DeleteVenue)
			})

			r.Route("/restore", func(r chi.Router) {
				r.Patch("/{id}", venuehandler.RestoreVenue)
			})
		})
	})

	return r
}
