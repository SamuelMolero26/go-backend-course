package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	db *pgxpool.Pool
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, email, password, role_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := s.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.RoleID,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil

}
