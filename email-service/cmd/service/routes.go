package main

import (
	"log"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors" // CORS IMPORT - REMOVE FOR PRODUCTION
	"github.com/orduro/pos-microservices/email-service/internal/handler/email"
	"github.com/orduro/pos-microservices/email-service/internal/handler/health"
)

func (s *server) mount() *chi.Mux {
	healthHandler := health.New(s.config.Env)
	emailHandler := email.New(s.store)

	r := chi.NewRouter()

	// ==========================================
	// CORS CONFIGURATION - FOR DEVELOPMENT ONLY, REMOVE THIS ENTIRE BLOCK FOR PRODUCTION

	if s.config.Env == "dev" || s.config.Env == "development" {
		corsOptions := cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}
		r.Use(cors.Handler(corsOptions))

		s.logCORSSettings(corsOptions)
	}

	// ==========================================

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

func (s *server) logCORSSettings(options cors.Options) {
	log.Printf("=== EMAIL SERVICE CORS CONFIGURATION ===")
	log.Printf("Environment: %s", s.config.Env)
	log.Printf("Allowed Origins: %v", options.AllowedOrigins)
	log.Printf("Allowed Methods: %v", options.AllowedMethods)
	log.Printf("Allowed Headers: %v", options.AllowedHeaders)
	log.Printf("Exposed Headers: %v", options.ExposedHeaders)
	log.Printf("Allow Credentials: %t", options.AllowCredentials)
	log.Printf("==========================================")
}
