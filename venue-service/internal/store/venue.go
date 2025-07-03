package store

import (
	"context"
	"database/sql"
	"time"
)

type Venue struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Address     string    `json:"address" validate:"required,min=1,max=500"`
	Phone       string    `json:"phone" validate:"omitempty,min=10,max=20"`
	VenueType   string    `json:"venue_type" validate:"omitempty,max=100"`
	Description string    `json:"description" validate:"required,min=1,max=1000"`
	Archived    bool      `json:"archived"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ValidateVenue(venue *Venue) error {
	return v.Struct(venue)
}

type VenueStore struct {
	db *sql.DB
}

func (s *VenueStore) Create(ctx context.Context, venue *Venue) error {
	query := `
	INSERT INTO venues (name, address, phone, venue_type, description)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`

	args := []any{venue.Name, venue.Address, venue.Phone, venue.VenueType, venue.Description}

	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&venue.ID,
		&venue.CreatedAt,
		&venue.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *VenueStore) GetByID(ctx context.Context, id int64) (*Venue, error) {
	query := `
	SELECT id, name, address, phone, venue_type, description, archived, created_at, updated_at
	FROM venues
	WHERE id = $1
	`

	v := Venue{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&v.ID,
		&v.Name,
		&v.Address,
		&v.Phone,
		&v.VenueType,
		&v.Description,
		&v.Archived,
		&v.CreatedAt,
		&v.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &v, nil
}
