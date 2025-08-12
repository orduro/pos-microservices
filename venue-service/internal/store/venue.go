package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/orduro/pos-microservices/common/auth"
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
	TenantID    string    `json:"tenant_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ValidateVenue(venue *Venue) error {
	return v.Struct(venue)
}

type VenueFilter struct {
	IncludeArchived bool
	Search          string
	Limit           int
	Offset          int
}

type VenueStore struct {
	db *sql.DB
}

func (s *VenueStore) Create(ctx context.Context, venue *Venue) error {
	// get tenant id from context
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return ErrTenantNotFound
	}

	query := `
	INSERT INTO venues (name, address, phone, venue_type, description, tenant_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, archived, created_at, updated_at
	`
	args := []any{venue.Name, venue.Address, venue.Phone, venue.VenueType, venue.Description, tenantID}
	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&venue.ID,
		&venue.Archived,
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
	venue.TenantID = *tenantID
	return nil
}

func (s *VenueStore) GetByID(ctx context.Context, id int64) (*Venue, error) {
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return nil, ErrTenantNotFound
	}

	query := `
	SELECT id, name, address, phone, venue_type, description, archived, tenant_id, created_at, updated_at
	FROM venues
	WHERE id = $1 AND tenant_id = $2
	`

	v := Venue{}

	args := []any{id, tenantID}

	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&v.ID,
		&v.Name,
		&v.Address,
		&v.Phone,
		&v.VenueType,
		&v.Description,
		&v.Archived,
		&v.TenantID,
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
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return ErrTenantNotFound
	}

	args := []any{id, tenantID}

	query := `UPDATE venues SET archived = true WHERE id = $1 AND tenant_id = $2 AND archived = false`
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// check if venue exists at all
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM venues WHERE id = $1 AND tenant_id = $2)`
		err := s.db.QueryRowContext(ctx, checkQuery, id, tenantID).Scan(&exists)
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
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return ErrTenantNotFound
	}

	query := `DELETE FROM venues WHERE id = $1 AND tenant_id = $2 AND archived = true`
	result, err := s.db.ExecContext(ctx, query, id, tenantID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// check if venue exists at all
		// or is archived
		var exists bool
		var archived bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM venues WHERE id = $1 AND tenant_id = $2), 
                      COALESCE((SELECT archived FROM venues WHERE id = $1 AND tenant_id = $2), false)`
		err := s.db.QueryRowContext(ctx, checkQuery, id, tenantID).Scan(&exists, &archived)
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

func (s *VenueStore) Update(ctx context.Context, venue *Venue) error {
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return ErrTenantNotFound
	}

	query := `
	UPDATE venues 
	SET name = $3, address = $4, phone = $5, venue_type = $6, description = $7, updated_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND tenant_id = $2 AND archived = false
	RETURNING updated_at
	`
	args := []any{venue.ID, tenantID, venue.Name, venue.Address, venue.Phone, venue.VenueType, venue.Description}

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&venue.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// check if venue exists at all
			var exists bool
			var archived bool

			checkQuery := `SELECT EXISTS(SELECT 1 FROM venues WHERE id = $1 AND tenant_id = $2), 
			              COALESCE((SELECT archived FROM venues WHERE id = $1 AND tenant_id = $2), false)`

			err := s.db.QueryRowContext(ctx, checkQuery, venue.ID, tenantID).Scan(&exists, &archived)

			if err != nil {
				return err
			}

			if !exists {
				return ErrVenueNotFound
			}
			if archived {
				return ErrVenueAlreadyArchived
			}
		}

		// check if its a unique constraint violation
		if database.IsPostgresUniqueConstraintError(err, "unique_venue_name_address") {
			return ErrVenueDuplicate
		}
		return err
	}

	venue.TenantID = *tenantID
	return nil
}

func (s *VenueStore) Restore(ctx context.Context, id int64) error {
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return ErrTenantNotFound
	}

	query := `UPDATE venues SET archived = false WHERE id = $1 AND tenant_id = $2 AND archived = true`
	result, err := s.db.ExecContext(ctx, query, id, tenantID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// check if venue exists at all
		var exists bool
		var archived bool

		checkQuery := `SELECT EXISTS(SELECT 1 FROM venues WHERE id = $1 AND tenant_id = $2), 
		              COALESCE((SELECT archived FROM venues WHERE id = $1 AND tenant_id = $2), false)`

		err := s.db.QueryRowContext(ctx, checkQuery, id, tenantID).Scan(&exists, &archived)

		if err != nil {
			return err
		}

		if !exists {
			return ErrVenueNotFound
		}
		if !archived {
			return ErrVenueNotArchived
		}
	}

	return nil
}

func (s *VenueStore) List(ctx context.Context, filter VenueFilter) ([]Venue, int, error) {
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return nil, 0, ErrTenantNotFound
	}

	// build the WHERE clause based on filters
	whereClause := "WHERE tenant_id = $1"
	args := []any{tenantID}
	argCount := 1

	// add archived filter
	if !filter.IncludeArchived {
		whereClause += " AND archived = false"
	}

	// add search filter
	// whatever is put in the search will check name address and description
	// to find matching
	if filter.Search != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR address ILIKE $%d OR description ILIKE $%d)",
			argCount, argCount, argCount)
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM venues %s", whereClause)
	var total int
	err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get venue count: %w", err)
	}

	// if no venues found, return early
	if total == 0 {
		return []Venue{}, 0, nil
	}

	// build the main query with pagination
	query := fmt.Sprintf(`
		SELECT id, name, address, phone, venue_type, description, archived, tenant_id, created_at, updated_at
		FROM venues
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount+1, argCount+2)

	// add limit and offset to args
	args = append(args, filter.Limit, filter.Offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query venues: %w", err)
	}
	defer rows.Close()

	var venues []Venue
	for rows.Next() {
		var venue Venue
		err := rows.Scan(
			&venue.ID,
			&venue.Name,
			&venue.Address,
			&venue.Phone,
			&venue.VenueType,
			&venue.Description,
			&venue.Archived,
			&venue.TenantID,
			&venue.CreatedAt,
			&venue.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan venue: %w", err)
		}
		venues = append(venues, venue)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating venue rows: %w", err)
	}

	return venues, total, nil
}
