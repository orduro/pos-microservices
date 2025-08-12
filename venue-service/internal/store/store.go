package store

import (
	"context"
	"database/sql"
)

type Store struct {
	Venues VenueRepository
	Items  ItemRepository
}

func NewStore(postgres *sql.DB) Store {
	return Store{
		Venues: &VenueStore{
			db: postgres,
		},
		Items: &ItemStore{
			db: postgres,
		},
	}
}

type VenueRepository interface {
	Create(ctx context.Context, venue *Venue) error
	GetByID(ctx context.Context, id int64) (*Venue, error)
	Archive(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, venue *Venue) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, filter VenueFilter) ([]Venue, int, error)
}

type ItemRepository interface {
	Create(ctx context.Context, item *Item) error
	GetByID(ctx context.Context, id int64) (*Item, error)
	Archive(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, item *Item) error
	Restore(ctx context.Context, id int64) error
	List(ctx context.Context, filter ItemFilter) ([]Item, int, error)
	ListByVenue(ctx context.Context, venueID int64, filter ItemFilter) ([]Item, int, error)
}
