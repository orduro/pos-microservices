package main

import (
	"log"

	"github.com/orduro/pos-microservices/auth-service/internal/database"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

func main() {
	// get config from .env
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// connect to postgres db
	pgconn, err := database.NewPostgresPool(cfg.Postgres.DSN, cfg.Postgres.MaxOpenConns, cfg.Postgres.MaxIdleConns, cfg.Postgres.MaxIdleTime)
	if err != nil {
		log.Fatal(err)
	}
	defer pgconn.Close()
	log.Println("postgres database connection pool established")

	store := store.NewStore(pgconn)

	srv := &server{
		config: cfg,
		store:  store,
	}

	// routes using chi router (routes.go)
	mux := srv.mount()

	// start server
	log.Fatal(srv.serve(mux))
}
