package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/orduro/pos-microservices/auth-service/internal/constants"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

type TokenService struct {
	token store.VerificationTokenRepository
}

func NewTokenService(store store.VerificationTokenRepository) *TokenService {
	return &TokenService{
		token: store,
	}
}

func (s *TokenService) GenerateToken() (string, error) {
	bytes := make([]byte, constants.TokenByteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *TokenService) CreateVerificationToken(ctx context.Context, userID int64) (string, error) {
	token, err := s.GenerateToken()
	if err != nil {
		return "", err
	}

	verificationToken := store.VerificationToken{
		UserID:    userID,
		Token:     token,
		TokenType: constants.TokenTypeEmailVerification,
		ExpiresAt: time.Now().Add(constants.TokenExpiryDuration),
	}

	if err := s.token.Create(ctx, &verificationToken); err != nil {
		return "", err
	}

	return token, nil
}

func (s *TokenService) CreatePasswordResetToken(ctx context.Context, userID int64) (string, error) {
	token, err := s.GenerateToken()
	if err != nil {
		return "", err
	}

	resetToken := store.VerificationToken{
		UserID:    userID,
		Token:     token,
		TokenType: constants.TokenTypePasswordReset,
		ExpiresAt: time.Now().Add(constants.TokenExpiryDuration),
	}

	if err := s.token.Create(ctx, &resetToken); err != nil {
		return "", err
	}

	return token, nil
}

func (s *TokenService) ValidateToken(ctx context.Context, tokenStr, tokenType string) (*store.VerificationToken, error) {
	token, err := s.token.GetByToken(ctx, tokenStr, tokenType)
	if err != nil {
		return nil, err
	}

	// check if token is expired
	if time.Now().After(token.ExpiresAt) {
		return nil, store.ErrTokenExpired
	}

	// mark token as used
	if err := s.token.MarkAsUsed(ctx, token.ID); err != nil {
		return nil, err
	}

	return token, nil
}
