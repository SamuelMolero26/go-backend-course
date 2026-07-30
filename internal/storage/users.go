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

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {

	query := `
		  SELECT id, username, email, password, role_id, created_at
		  FROM users WHERE id = $1`

	var user User
	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.Password, &user.RoleID, &user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (s *UserStore) Update(ctx context.Context, user *User) error {
	query := `UPDATE users SET username = $1, email = $2 WHERE id = $3`

	result, err := s.db.Exec(ctx, query, user.Username, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password, role_id, created_at
		FROM users WHERE email = $1`

	var user User
	err := s.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.Password, &user.RoleID, &user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}
