package store

import (
	"context"
	"database/sql"
)

type Store struct {
	Venues VenueRepository
}

func NewStore(postgres *sql.DB) Store {
	return Store{
		Venues: &VenueStore{
			db: postgres,
		},
	}
}

type VenueRepository interface {
	Create(ctx context.Context, venue *Venue) error
	GetByID(ctx context.Context, id int64) (*Venue, error)
}
