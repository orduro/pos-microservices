package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/orduro/pos-microservices/common/auth"
	"github.com/orduro/pos-microservices/venue/internal/database"
)

type Modifier struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Required    bool    `json:"required"`
	MaxQuantity int     `json:"max_quantity,omitempty"`
}

type Item struct {
	ID          int64      `json:"id"`
	VenueID     int64      `json:"venue_id" validate:"required"`
	Name        string     `json:"name" validate:"required,min=1"`
	Description string     `json:"description" validate:"required,min=1"`
	Category    string     `json:"category" validate:"required,min=1"`
	Price       float64    `json:"price" validate:"required,min=0"`
	IsAvailable bool       `json:"is_available"`
	ImageURL    string     `json:"image_url" validate:"omitempty,url"`
	Position    int        `json:"position" validate:"min=0"`
	Modifiers   []Modifier `json:"modifiers"`
	Archived    bool       `json:"archived"`
	TenantID    string     `json:"tenant_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ValidateItem(item *Item) error {
	return v.Struct(item)
}

type ItemFilter struct {
	VenueID         int64
	IncludeArchived bool
	OnlyAvailable   bool
	Category        string
	Search          string
	Limit           int
	Offset          int
}

type ItemStore struct {
	db *sql.DB
}

func (s *ItemStore) Create(ctx context.Context, item *Item) error {
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return ErrTenantNotFound
	}

	modifiersJSON, err := json.Marshal(item.Modifiers)
	if err != nil {
		return fmt.Errorf("failed to marshal modifiers: %w", err)
	}

	query := `
	INSERT INTO items (venue_id, name, description, category, price, is_available, image_url, position, modifiers, tenant_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id, archived, created_at, updated_at
	`
	args := []any{item.VenueID, item.Name, item.Description, item.Category, item.Price, item.IsAvailable, item.ImageURL, item.Position, modifiersJSON, tenantID}
	err = s.db.QueryRowContext(ctx, query, args...).Scan(
		&item.ID,
		&item.Archived,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if database.IsPostgresUniqueConstraintError(err, "items_venue_id_name_tenant_id_key") {
			return ErrItemDuplicate
		}
		return err
	}
	item.TenantID = *tenantID
	return nil
}

func (s *ItemStore) GetByID(ctx context.Context, id int64) (*Item, error) {
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return nil, ErrTenantNotFound
	}

	query := `
	SELECT id, venue_id, name, description, category, price, is_available, image_url, position, modifiers, archived, tenant_id, created_at, updated_at
	FROM items
	WHERE id = $1 AND tenant_id = $2
	`

	item := Item{}
	var modifiersJSON []byte

	err := s.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&item.ID,
		&item.VenueID,
		&item.Name,
		&item.Description,
		&item.Category,
		&item.Price,
		&item.IsAvailable,
		&item.ImageURL,
		&item.Position,
		&modifiersJSON,
		&item.Archived,
		&item.TenantID,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrItemNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(modifiersJSON, &item.Modifiers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal modifiers: %w", err)
	}

	return &item, nil
}

func (s *ItemStore) Archive(ctx context.Context, id int64) error {
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return ErrTenantNotFound
	}

	query := `UPDATE items SET archived = true WHERE id = $1 AND tenant_id = $2 AND archived = false`
	result, err := s.db.ExecContext(ctx, query, id, tenantID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM items WHERE id = $1 AND tenant_id = $2)`
		err := s.db.QueryRowContext(ctx, checkQuery, id, tenantID).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			return ErrItemNotFound
		}
		return ErrItemAlreadyArchived
	}

	return nil
}

func (s *ItemStore) Delete(ctx context.Context, id int64) error {
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return ErrTenantNotFound
	}

	query := `DELETE FROM items WHERE id = $1 AND tenant_id = $2 AND archived = true`
	result, err := s.db.ExecContext(ctx, query, id, tenantID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		var exists bool
		var archived bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM items WHERE id = $1 AND tenant_id = $2), 
                      COALESCE((SELECT archived FROM items WHERE id = $1 AND tenant_id = $2), false)`
		err := s.db.QueryRowContext(ctx, checkQuery, id, tenantID).Scan(&exists, &archived)
		if err != nil {
			return err
		}

		if !exists {
			return ErrItemNotFound
		}
		return ErrItemNotArchived
	}

	return nil
}

func (s *ItemStore) Update(ctx context.Context, item *Item) error {
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return ErrTenantNotFound
	}

	modifiersJSON, err := json.Marshal(item.Modifiers)
	if err != nil {
		return fmt.Errorf("failed to marshal modifiers: %w", err)
	}

	query := `
	UPDATE items 
	SET venue_id = $3, name = $4, description = $5, category = $6, price = $7, is_available = $8, image_url = $9, position = $10, modifiers = $11, updated_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND tenant_id = $2 AND archived = false
	RETURNING updated_at
	`
	args := []any{item.ID, tenantID, item.VenueID, item.Name, item.Description, item.Category, item.Price, item.IsAvailable, item.ImageURL, item.Position, modifiersJSON}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&item.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			var exists bool
			var archived bool

			checkQuery := `SELECT EXISTS(SELECT 1 FROM items WHERE id = $1 AND tenant_id = $2), 
			              COALESCE((SELECT archived FROM items WHERE id = $1 AND tenant_id = $2), false)`

			err := s.db.QueryRowContext(ctx, checkQuery, item.ID, tenantID).Scan(&exists, &archived)
			if err != nil {
				return err
			}

			if !exists {
				return ErrItemNotFound
			}
			if archived {
				return ErrItemAlreadyArchived
			}
		}

		if database.IsPostgresUniqueConstraintError(err, "items_venue_id_name_tenant_id_key") {
			return ErrItemDuplicate
		}
		return err
	}

	item.TenantID = *tenantID
	return nil
}

func (s *ItemStore) Restore(ctx context.Context, id int64) error {
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return ErrTenantNotFound
	}

	query := `UPDATE items SET archived = false WHERE id = $1 AND tenant_id = $2 AND archived = true`
	result, err := s.db.ExecContext(ctx, query, id, tenantID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		var exists bool
		var archived bool

		checkQuery := `SELECT EXISTS(SELECT 1 FROM items WHERE id = $1 AND tenant_id = $2), 
		              COALESCE((SELECT archived FROM items WHERE id = $1 AND tenant_id = $2), false)`

		err := s.db.QueryRowContext(ctx, checkQuery, id, tenantID).Scan(&exists, &archived)
		if err != nil {
			return err
		}

		if !exists {
			return ErrItemNotFound
		}
		if !archived {
			return ErrItemNotArchived
		}
	}

	return nil
}

func (s *ItemStore) List(ctx context.Context, filter ItemFilter) ([]Item, int, error) {
	tenantID, ok := auth.GetTenantID(ctx)
	if !ok || tenantID == nil || *tenantID == "" {
		return nil, 0, ErrTenantNotFound
	}

	whereClause := "WHERE tenant_id = $1"
	args := []any{tenantID}
	argCount := 1

	if filter.VenueID > 0 {
		argCount++
		whereClause += fmt.Sprintf(" AND venue_id = $%d", argCount)
		args = append(args, filter.VenueID)
	}

	if !filter.IncludeArchived {
		whereClause += " AND archived = false"
	}

	if filter.OnlyAvailable {
		whereClause += " AND is_available = true"
	}

	if filter.Category != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, filter.Category)
	}

	if filter.Search != "" {
		argCount++
		whereClause += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM items %s", whereClause)
	var total int
	err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get item count: %w", err)
	}

	if total == 0 {
		return []Item{}, 0, nil
	}

	query := fmt.Sprintf(`
		SELECT id, venue_id, name, description, category, price, is_available, image_url, position, modifiers, archived, tenant_id, created_at, updated_at
		FROM items
		%s
		ORDER BY position ASC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount+1, argCount+2)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		var modifiersJSON []byte
		err := rows.Scan(
			&item.ID,
			&item.VenueID,
			&item.Name,
			&item.Description,
			&item.Category,
			&item.Price,
			&item.IsAvailable,
			&item.ImageURL,
			&item.Position,
			&modifiersJSON,
			&item.Archived,
			&item.TenantID,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan item: %w", err)
		}

		if err := json.Unmarshal(modifiersJSON, &item.Modifiers); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal modifiers: %w", err)
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating item rows: %w", err)
	}

	return items, total, nil
}

func (s *ItemStore) ListByVenue(ctx context.Context, venueID int64, filter ItemFilter) ([]Item, int, error) {
	filter.VenueID = venueID
	return s.List(ctx, filter)
}
