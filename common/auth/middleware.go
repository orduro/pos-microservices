package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/orduro/pos-microservices/common/json"
)

func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				json.WriteError(w, r, http.StatusUnauthorized, "missing token")
				return
			}

			claims, err := ValidateJWTToken(token, secret)
			if err != nil {
				json.WriteError(w, r, http.StatusUnauthorized, err.Error())
				log.Println(err)
				return
			}

			// don't allow refresh tokens for API access
			if claims.IsRefresh {
				json.WriteError(w, r, http.StatusUnauthorized, "refresh tokens cannot be used for API access")
				return
			}

			userInfo := UserInfo{
				UserID:   claims.UserID,
				Email:    claims.Email,
				TenantID: claims.TenantID,
			}

			ctx := context.WithValue(r.Context(), userContextKey, userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
