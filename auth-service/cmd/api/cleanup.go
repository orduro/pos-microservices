package main

import (
	"context"
	"log"
	"time"

	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

func startTokenCleanupJob(tokenStore store.VerificationTokenRepository) {
	// run once per day
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	log.Println("started token cleanup job - runs every 24 hours")

	// run cleanup immediately on startup
	if err := cleanupExpiredTokens(tokenStore); err != nil {
		log.Printf("initial token cleanup failed: %v", err)
	}

	// then run on schedule
	for {
		select {
		case <-ticker.C:
			if err := cleanupExpiredTokens(tokenStore); err != nil {
				log.Printf("scheduled token cleanup failed: %v", err)
			}
		}
	}
}

func cleanupExpiredTokens(tokenStore store.VerificationTokenRepository) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("running expired token cleanup...")

	if err := tokenStore.DeleteExpired(ctx); err != nil {
		return err
	}

	log.Println("expired token cleanup completed successfully")
	return nil
}
