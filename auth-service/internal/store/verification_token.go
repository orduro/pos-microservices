package store

import (
	"context"
	"database/sql"
	"time"
)

type Token struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Token     string     `json:"token"`
	TokenType string     `json:"token_type"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type VerificationTokenStore struct {
	db *sql.DB
}

func (s *VerificationTokenStore) Create(ctx context.Context, token *Token) error {
	query := `
		INSERT INTO verification_tokens (user_id, token, token_type, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	args := []any{token.UserID, token.Token, token.TokenType, token.ExpiresAt}

	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&token.ID,
		&token.CreatedAt,
		&token.UpdatedAt,
	)
	return err
}

func (s *VerificationTokenStore) GetByToken(ctx context.Context, token string, tokenType string) (*Token, error) {
	query := `
		SELECT id, user_id, token, token_type, expires_at, used_at, created_at, updated_at
		FROM verification_tokens
		WHERE token = $1 AND token_type = $2 AND used_at IS NULL
	`

	vt := &Token{}
	err := s.db.QueryRowContext(ctx, query, token, tokenType).Scan(
		&vt.ID,
		&vt.UserID,
		&vt.Token,
		&vt.TokenType,
		&vt.ExpiresAt,
		&vt.UsedAt,
		&vt.CreatedAt,
		&vt.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	return vt, nil
}

func (s *VerificationTokenStore) MarkAsUsed(ctx context.Context, tokenID int64) error {
	query := `
		UPDATE verification_tokens 
		SET used_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND used_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTokenNotFound
	}

	return nil
}

// clean up expired tokens - this will be called periodically (once per day?)
func (s *VerificationTokenStore) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM verification_tokens WHERE expires_at < NOW()`
	_, err := s.db.ExecContext(ctx, query)
	return err
}
