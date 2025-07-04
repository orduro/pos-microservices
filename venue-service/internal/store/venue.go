package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/orduro/pos-microservices/venue/internal/database"
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
		// Check if it's a unique constraint violation
		if database.IsPostgresUniqueConstraintError(err, "unique_venue_name_address") {
			return ErrVenueDuplicate
		}
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
		if err == sql.ErrNoRows {
			return nil, ErrVenueNotFound
		}
		return nil, err
	}

	return &v, nil
}

func (s *VenueStore) Archive(ctx context.Context, id int64) error {
	query := `UPDATE venues SET archived = true WHERE id = $1 AND archived = false`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Check if venue exists at all
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM venues WHERE id = $1)`
		err := s.db.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			return ErrVenueNotFound
		}
		return ErrVenueAlreadyArchived
	}

	return nil
}

func (s *VenueStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM venues WHERE id = $1 AND archived = true`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// Check if venue exists at all
		// Or is archived
		var exists bool
		var archived bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM venues WHERE id = $1), 
                      COALESCE((SELECT archived FROM venues WHERE id = $1), false)`
		err := s.db.QueryRowContext(ctx, checkQuery, id).Scan(&exists, &archived)
		if err != nil {
			return err
		}

		if !exists {
			return ErrVenueNotFound
		}
		return ErrVenueNotArchived
	}

	return nil
}
