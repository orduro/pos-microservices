package service

import (
	"time"

	"github.com/orduro/common/auth"
)

// JWTService wraps the common JWT service for auth-service specific needs
type JWTService struct {
	*auth.JWTService
}

// TokenPair alias for consistency with existing code
type TokenPair = auth.TokenPair

func NewJWTService(secret string, expirationTime, refreshExpirationTime time.Duration) *JWTService {
	return &JWTService{
		JWTService: auth.NewJWTService(secret, expirationTime, refreshExpirationTime),
	}
}

