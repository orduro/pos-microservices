// Update your auth-service/internal/store/store.go
package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
)

var (
	V                = validator.New()
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrTokenNotFound = errors.New("verification token not found")
	ErrTokenExpired  = errors.New("verification token expired")
)

type Store struct {
	Users              UserRepository
	VerificationTokens VerificationTokenRepository
}

func NewStore(postgres *sql.DB) Store {
	return Store{
		Users: &UserStore{
			db: postgres,
		},
		VerificationTokens: &VerificationTokenStore{
			db: postgres,
		},
	}
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	MarkAsVerified(ctx context.Context, userID int64) error
}

type VerificationTokenRepository interface {
	Create(ctx context.Context, token *VerificationToken) error
	GetByToken(ctx context.Context, token string, tokenType string) (*VerificationToken, error)
	MarkAsUsed(ctx context.Context, tokenID int64) error
	DeleteExpired(ctx context.Context) error
}
