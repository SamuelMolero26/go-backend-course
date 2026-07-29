package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FollowerStore struct {
	db *pgxpool.Pool
}

func (s *FollowerStore) Follow(ctx context.Context, userID, followerID int64) error {
	query := `INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING`

	_, err := s.db.Exec(ctx, query, userID, followerID)
	if err != nil {
		return fmt.Errorf("follow: %w", err)
	}
	return nil
}

func (s *FollowerStore) Unfollow(ctx context.Context, userID, followerID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`

	result, err := s.db.Exec(ctx, query, userID, followerID)
	if err != nil {
		return fmt.Errorf("unfollow: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("follow relationship not found")
	}

	return nil
}

func (s *FollowerStore) GetFollowers(ctx context.Context, userID int64) ([]User, error) {
	query := `SELECT u.id, u.username, u.email, u.created_at
		FROM followers f
		JOIN users u ON u.id = f.follower_id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get followers: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan follower: %w", err)
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration: %w", rows.Err())
	}

	return users, nil
}

func (s *FollowerStore) GetFollowing(ctx context.Context, userID int64) ([]User, error) {
	query := `SELECT u.id, u.username, u.email, u.created_at
		FROM followers f
		JOIN users u ON u.id = f.user_id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get following: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan following: %w", err)
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration: %w", rows.Err())
	}

	return users, nil
}
