package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentStore struct {
	db *pgxpool.Pool
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := s.db.QueryRow(ctx, query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt)

	if err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `SELECT id, post_id, user_id, content, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC`

	rows, err := s.db.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("get comments by post: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.UserID,
			&comment.Content, &comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration: %w", rows.Err())
	}

	return comments, nil
}

func (s *CommentStore) GetByID(ctx context.Context, id int64) (*Comment, error) {
	query := `SELECT id, post_id, user_id, content, created_at FROM comments WHERE id = $1`

	var comment Comment
	err := s.db.QueryRow(ctx, query, id).Scan(
		&comment.ID, &comment.PostID, &comment.UserID,
		&comment.Content, &comment.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("comment not found: %w", err)
		}
		return nil, fmt.Errorf("get comment by id: %w", err)
	}
	return &comment, nil
}

func (s *CommentStore) Delete(ctx context.Context, id int64, userID int64) error {
	query := `DELETE FROM comments WHERE id = $1 AND user_id = $2`

	result, err := s.db.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment not found or not owned by user")
	}

	return nil
}
