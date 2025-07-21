package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret                []byte
	expirationTime        time.Duration
	refreshExpirationTime time.Duration
}

type Claims struct {
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	IsRefresh bool   `json:"is_refresh,omitempty"`
	jwt.RegisteredClaims
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

func (s *JWTService) GenerateTokenPair(userID int64, email string) (*TokenPair, error) {
	// generate access token
	accessToken, accessExpiresAt, err := s.generateToken(userID, email, false, s.expirationTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, _, err := s.generateToken(userID, email, true, s.refreshExpirationTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiresAt.Unix(),
	}, nil
}

func (s *JWTService) generateToken(userID int64, email string, isRefresh bool, duration time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().Add(duration)

	claims := Claims{
		UserID:    userID,
		Email:     email,
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

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTService) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if !claims.IsRefresh {
		return nil, errors.New("provided token is not a refresh token")
	}

	// generate new token pair
	return s.GenerateTokenPair(claims.UserID, claims.Email)
}
