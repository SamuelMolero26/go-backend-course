package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostStore struct {
	db *pgxpool.Pool
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (user_id, content)
		VALUES ($1, $2)
		RETURNING id, created_at`

	err := s.db.QueryRow(ctx, query,
		post.UserID,
		post.Content,
	).Scan(&post.ID, &post.CreatedAt)

	if err != nil {
		return fmt.Errorf("create post: %w", err)
	}
	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id, user_id, content, created_at
		FROM posts WHERE id = $1`

	var post Post
	err := s.db.QueryRow(ctx, query, id).Scan(
		&post.ID, &post.UserID, &post.Content, &post.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("post not found: %w", err)
		}
		return nil, fmt.Errorf("get post by id: %w", err)
	}
	return &post, nil
}

func (s *PostStore) GetByUserID(ctx context.Context, userID int64) ([]Post, error) {
	query := `SELECT id, user_id, content, created_at
		FROM posts
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get posts by user: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration: %w", rows.Err())
	}

	return posts, nil
}

func (s *PostStore) Delete(ctx context.Context, id int64, userID int64) error {
	query := `DELETE FROM posts WHERE id = $1 AND user_id = $2`

	result, err := s.db.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("post not found or not owned by user")
	}

	return nil
}
