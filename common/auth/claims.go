package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID    int64   `json:"user_id"`
	Email     string  `json:"email"`
	TenantID  *string `json:"tenant_id,omitempty"` // for multitenancy
	IsRefresh bool    `json:"is_refresh,omitempty"`
	jwt.RegisteredClaims
}

type UserInfo struct {
	UserID   int64
	Email    string
	TenantID *string
}
