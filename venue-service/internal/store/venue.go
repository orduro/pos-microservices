package store

import (
	"context"
	"database/sql"
	"time"
)

type Venue struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	VenueType   string    `json:"venue_type"`
	Description string    `json:"description"`
	Archived    bool      `json:"archived"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type VenueStore struct {
	db *sql.DB
}

func (s *VenueStore) Create(ctx context.Context, venue *Venue) error {
	return nil
}
