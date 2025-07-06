package store

import (
	"context"
	"database/sql"
)

type Store struct {
	Users UserRepository
}

func NewStore(postgres *sql.DB) Store {
	return Store{
		Users: &UserStore{
			db: postgres,
		},
	}
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id int64) (*User, error)
}
