package main

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/orduro/pos-microservices/communication-service/internal/handler/email"
	"github.com/orduro/pos-microservices/communication-service/internal/handler/health"
)

func (s *server) mount() *chi.Mux {
	healthHandler := health.New(s.config.Env)
	emailHandler := email.New(s.store)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", healthHandler.Check)

		// routes for /email
		r.Route("/email", func(r chi.Router) {
			r.Post("/send", emailHandler.SendEmail)
			r.Post("/verification", emailHandler.SendVerificationEmail)
			r.Post("/password-reset", emailHandler.SendPasswordResetEmail)
			r.Post("/welcome", emailHandler.SendWelcomeEmail)
		})
	})

	return r
}
