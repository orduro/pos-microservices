package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/orduro/pos-microservices/auth-service/internal/database"
)

type UserRegistrationDetails struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

func (u *UserRegistrationDetails) Validate() error {
	return V.Struct(u)
}

type UserLoginDetails struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (u *UserLoginDetails) Validate() error {
	return V.Struct(u)
}

type User struct {
	ID           int64      `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	IsVerified   bool       `json:"is_verified"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	TenantID     string     `json:"tenant_id"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, is_verified, tenant_id, created_at, updated_at
	`
	args := []any{user.Email, user.PasswordHash, user.FirstName, user.LastName}

	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.IsVerified,
		&user.TenantID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		// Check if it's a unique constraint violation for email
		if database.IsPostgresUniqueConstraintError(err, "email") {
			return ErrUserExists
		}
		return err
	}

	return nil
}
func (s *UserStore) GetById(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, last_login_at, is_verified, tenant_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.LastLoginAt,
		&user.IsVerified,
		&user.TenantID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, last_login_at, is_verified, tenant_id, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.LastLoginAt,
		&user.IsVerified,
		&user.TenantID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserStore) MarkAsVerified(ctx context.Context, userID int64) error {
	query := `
		UPDATE users 
		SET is_verified = true, updated_at = NOW()
		WHERE id = $1 AND is_verified = false
	`

	result, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (s *UserStore) UpdateLastLogin(ctx context.Context, userID int64) error {
	query := `
		UPDATE users 
		SET last_login_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`

	result, err := s.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (s *UserStore) UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error {
	query := `
	UPDATE users
	SET password_hash = $1, updated_at = NOW()
	WHERE id = $2
	`

	args := []any{hashedPassword, userID}

	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
