package store

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	LastLoginAt  time.Time
	IsVerified   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	return nil
}
func (s *UserStore) GetById(ctx context.Context, id int64) (*User, error) {
	return nil, nil
}
