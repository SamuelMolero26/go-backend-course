package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleStore struct {
	db *pgxpool.Pool
}

func (s *RoleStore) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `SELECT id, name FROM roles WHERE name = $1`

	var role Role
	err := s.db.QueryRow(ctx, query, name).Scan(&role.ID, &role.Name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("role not found: %w", err)
		}
		return nil, fmt.Errorf("get role by name: %w", err)
	}
	return &role, nil
}
