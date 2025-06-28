package main

import (
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/orduro/pos-microservices/venue/internal/handler/health"
)

func (app *application) mount() *chi.Mux {
	// initialise handlers
	healthhandler := health.New(app.config.Env)
	r := chi.NewRouter()

	// only add cors in development mode
	if app.config.Env == "dev" || app.config.Env == "development" {

		var allowedOrigins []string

		switch app.config.Env {
		case "dev", "development":
			allowedOrigins = []string{"*"}

		case "staging", "prod", "production":
			originsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")

			if originsEnv != "" {
				allowedOrigins = strings.Split(originsEnv, ",")
			}

		default:
			allowedOrigins = []string{"http://localhost:3000"}
		}

		// Basic CORS
		// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: allowedOrigins,
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// prefix /api in front of all routes, all routes go in here
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", healthhandler.Check)
	})

	return r
}
