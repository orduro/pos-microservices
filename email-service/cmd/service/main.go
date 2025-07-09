package main

import (
	"log"

	"github.com/orduro/pos-microservices/email-service/internal/store"
)

func main() {
	// get config from environment
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// initialise email store (no database needed for communications)
	emailConfig := store.EmailConfig{
		FromEmail:    cfg.Email.FromEmail,
		FromName:     cfg.Email.FromName,
		ReplyToEmail: cfg.Email.ReplyTo,
		Provider:     cfg.Email.Provider,
		SMTPHost:     cfg.Email.SMTPHost,
		SMTPPort:     cfg.Email.SMTPPort,
		SMTPUser:     cfg.Email.SMTPUser,
		SMTPPass:     cfg.Email.SMTPPass,
	}

	store := store.NewStore(emailConfig)

	srv := &server{
		config: cfg,
		store:  store,
	}

	// routes using chi router (routes.go)
	mux := srv.mount()

	// start server
	log.Fatal(srv.serve(mux))
}
