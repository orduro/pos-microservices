package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/orduro/pos-microservices/auth-service/internal/constants"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

type TokenService struct {
	store store.VerificationTokenRepository
}

func NewTokenService(store store.VerificationTokenRepository) *TokenService {
	return &TokenService{
		store: store,
	}
}

func (s *TokenService) GenerateToken() (string, error) {
	bytes := make([]byte, constants.TokenByteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func (s *TokenService) CreateVerificationToken(ctx context.Context, userID int64) (string, error) {
	token, err := s.GenerateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	verificationToken := store.VerificationToken{
		UserID:    userID,
		Token:     token,
		TokenType: constants.TokenTypeEmailVerification,
		ExpiresAt: time.Now().Add(constants.TokenExpiryDuration),
	}

	if err := s.store.Create(ctx, &verificationToken); err != nil {
		return "", fmt.Errorf("failed to store verification token: %w", err)
	}

	return token, nil
}

func (s *TokenService) CreatePasswordResetToken(ctx context.Context, userID int64) (string, error) {
	token, err := s.GenerateToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	resetToken := store.VerificationToken{
		UserID:    userID,
		Token:     token,
		TokenType: constants.TokenTypePasswordReset,
		ExpiresAt: time.Now().Add(constants.TokenExpiryDuration),
	}

	if err := s.store.Create(ctx, &resetToken); err != nil {
		return "", fmt.Errorf("failed to store password reset token: %w", err)
	}

	return token, nil
}

func (s *TokenService) ValidateToken(ctx context.Context, tokenStr, tokenType string) (*store.VerificationToken, error) {
	token, err := s.store.GetByToken(ctx, tokenStr, tokenType)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// check if token is expired
	if time.Now().After(token.ExpiresAt) {
		return nil, store.ErrTokenExpired
	}

	// mark token as used
	if err := s.store.MarkAsUsed(ctx, token.ID); err != nil {
		return nil, fmt.Errorf("failed to mark token as used: %w", err)
	}

	return token, nil
}
