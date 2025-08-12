package main

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/orduro/common/auth"
	"github.com/orduro/pos-microservices/venue/internal/handler/health"
	"github.com/orduro/pos-microservices/venue/internal/handler/item"
	"github.com/orduro/pos-microservices/venue/internal/handler/venue"
)

func (s *server) mount() *chi.Mux {
	// initialise handlers
	healthHandler := health.New(s.config.Env)
	venueHandler := venue.New(s.store)
	itemHandler := item.New(s.store)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", healthHandler.Check)

	// prefix /api in front of all routes, all routes go in here
	r.Route("/api", func(r chi.Router) {
		// all routes in /api require the user to be authenticated
		r.Use(auth.JWTAuth(s.config.JWT.Secret))

		// routes for /venue
		r.Route("/venue", func(r chi.Router) {
			r.Post("/", venueHandler.CreateVenue)
			r.Get("/", venueHandler.ListVenues)
			r.Patch("/{id}", venueHandler.UpdateVenue)

			// routes for venue archives `/venue/archive`
			r.Route("/archive", func(r chi.Router) {
				r.Put("/{id}", venueHandler.ArchiveVenue)
			})

			// routes for venue deletes `/venue/delete`
			r.Route("/delete", func(r chi.Router) {
				r.Delete("/{id}", venueHandler.DeleteVenue)
			})

			r.Route("/restore", func(r chi.Router) {
				r.Patch("/{id}", venueHandler.RestoreVenue)
			})
		})

		// routes for /item
		r.Route("/items", func(r chi.Router) {
			r.Post("/", itemHandler.CreateItem)
			r.Get("/", itemHandler.ListItems)
			r.Get("/{id}", itemHandler.GetItem)
			r.Patch("/{id}", itemHandler.UpdateItem)

			r.Route("/archive", func(r chi.Router) {
				r.Put("/{id}", itemHandler.ArchiveItem)
			})

			r.Route("/delete", func(r chi.Router) {
				r.Delete("/{id}", itemHandler.DeleteItem)
			})

			r.Route("/restore", func(r chi.Router) {
				r.Patch("/{id}", itemHandler.RestoreItem)
			})
		})
	})

	return r
}
