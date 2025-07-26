package main

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/orduro/pos-microservices/auth-service/internal/handler"
)

func (s *server) mount() *chi.Mux {
	// initialise handler
	handler := handler.New(s.config.Env, s.userService)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// healthcheck
	r.Get("/health", handler.Healthcheck)

	// prefix /api in front of all routes, all routes go in here
	r.Route("/api", func(r chi.Router) {

		r.Post("/register", handler.RegisterUser)
		r.Post("/login", handler.Login)
		r.Post("/refresh", handler.RefreshToken)
		r.Post("/resend-verification", handler.ResendVerificationEmail)
		r.Get("/verify/{token}", handler.UserVerificationCallback)
		r.Post("/forgot-password", handler.InitialiseForgotPassword)
	})

	return r
}
