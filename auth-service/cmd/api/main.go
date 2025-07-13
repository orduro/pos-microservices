package main

import (
	"log"

	"github.com/orduro/pos-microservices/auth-service/internal/database"
	"github.com/orduro/pos-microservices/auth-service/internal/httpclient"
	"github.com/orduro/pos-microservices/auth-service/internal/service"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

func main() {
	// get config from environment
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// connect to postgres db
	pgconn, err := database.NewPostgresPool(
		cfg.Postgres.DSN,
		cfg.Postgres.MaxOpenConns,
		cfg.Postgres.MaxIdleConns,
		cfg.Postgres.MaxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer pgconn.Close()
	log.Println("postgres database connection pool established")

	// init store
	store := store.NewStore(pgconn)

	// init httpclient for service-service communication
	httpClientConfig := httpclient.Config{
		CommunicationServiceURL: cfg.Services.CommunicationServiceURL,
	}
	httpClient := httpclient.NewHttpClient(httpClientConfig)

	// init services
	tokenService := service.NewTokenService(store.VerificationTokens)
	userService := service.NewUserService(store.Users, tokenService, httpClient, cfg.Services.FrontendAdminURL)

	srv := &server{
		config:      cfg,
		userService: userService,
	}

	// routes using chi router (routes.go)
	mux := srv.mount()

	// start server
	log.Fatal(srv.serve(mux))
}
