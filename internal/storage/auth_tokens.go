package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthTokenStore struct {
	db *pgxpool.Pool
}

func (s *AuthTokenStore) Create(ctx context.Context, token *AuthToken) error {
	query := `INSERT INTO auth_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := s.db.QueryRow(ctx, query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)

	if err != nil {
		return fmt.Errorf("create auth token: %w", err)
	}
	return nil
}

func (s *AuthTokenStore) GetByToken(ctx context.Context, token string) (*AuthToken, error) {
	query := `SELECT id, user_id, token, expires_at, created_at
		FROM auth_tokens
		WHERE token = $1 AND expires_at > NOW()`

	var authToken AuthToken
	err := s.db.QueryRow(ctx, query, token).Scan(
		&authToken.ID, &authToken.UserID, &authToken.Token,
		&authToken.ExpiresAt, &authToken.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("token not found or expired: %w", err)
		}
		return nil, fmt.Errorf("get token: %w", err)
	}
	return &authToken, nil
}

func (s *AuthTokenStore) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM auth_tokens WHERE token = $1`

	_, err := s.db.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("delete token: %w", err)
	}
	return nil
}

func (s *AuthTokenStore) DeleteByUserID(ctx context.Context, userID int64) error {
	query := `DELETE FROM auth_tokens WHERE user_id = $1`

	_, err := s.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("delete tokens by user: %w", err)
	}
	return nil
}
