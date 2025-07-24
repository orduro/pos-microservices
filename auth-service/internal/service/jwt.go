package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/orduro/common/auth"
)

type JWTService struct {
	secret                []byte
	expirationTime        time.Duration
	refreshExpirationTime time.Duration
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func NewJWTService(secret string, expirationTime, refreshExpirationTime time.Duration) *JWTService {
	return &JWTService{
		secret:                []byte(secret),
		expirationTime:        expirationTime,
		refreshExpirationTime: refreshExpirationTime,
	}
}

func (s *JWTService) GenerateTokenPairWithTenant(userID int64, email string, tenantID *string) (*TokenPair, error) {
	// generate access token
	accessToken, accessExpiresAt, err := s.generateToken(userID, email, tenantID, false, s.expirationTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// generate refresh token
	refreshToken, _, err := s.generateToken(userID, email, tenantID, true, s.refreshExpirationTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}

func (s *JWTService) generateToken(userID int64, email string, tenantID *string, isRefresh bool, duration time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().Add(duration)

	claims := auth.Claims{
		UserID:    userID,
		Email:     email,
		TenantID:  tenantID,
		IsRefresh: isRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (s *JWTService) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	// use validation function from common package
	claims, err := auth.ValidateJWTToken(refreshTokenString, string(s.secret))
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if !claims.IsRefresh {
		return nil, errors.New("provided token is not a refresh token")
	}

	// generate new token pair
	return s.GenerateTokenPairWithTenant(claims.UserID, claims.Email, claims.TenantID)
}
