package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store struct {
	db *pgxpool.Pool
	*UserStore
	*PostStore
	*CommentStore
	*FollowerStore
	*AuthTokenStore
	*RoleStore
}

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	RoleID    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Post struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty"`
}

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty"`
}

type Follower struct {
	UserID     int64     `json:"user_id"`
	FollowerID int64     `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type AuthToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type RoleRepository interface {
	GetByName(ctx context.Context, name string) (*Role, error)
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id int64) (*Post, error)
	GetByUserID(ctx context.Context, userID int64) ([]Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int64, userID int64) error
}

type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	GetByID(ctx context.Context, id int64) (*Comment, error) 
	Delete(ctx context.Context, id int64, userID int64) error
}

type FollowerRepository interface {
	Follow(ctx context.Context, userID, followerID int64) error
	Unfollow(ctx context.Context, userID, followerID int64) error
	GetFollowers(ctx context.Context, userID int64) ([]User, error)
	GetFollowing(ctx context.Context, userID int64) ([]User, error)
}

type AuthTokenRepository interface {
	Create(ctx context.Context, token *AuthToken) error
	GetByToken(ctx context.Context, token string) (*AuthToken, error)
	Delete(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID int64) error
}

type UsersRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type Storage interface {
	Posts() PostRepository
	Comments() CommentRepository
	Users() UsersRepository
	Followers() FollowerRepository
	AuthTokens() AuthTokenRepository
	Roles() RoleRepository
}



func (s *Store) Users() UsersRepository          { return s.UserStore }
func (s *Store) Posts() PostRepository           { return s.PostStore }
func (s *Store) Comments() CommentRepository     { return s.CommentStore }
func (s *Store) Followers() FollowerRepository   { return s.FollowerStore }
func (s *Store) AuthTokens() AuthTokenRepository { return s.AuthTokenStore }
func (s *Store) Roles() RoleRepository           { return s.RoleStore }

func (s *Store) Close() {
	s.db.Close()
}

func NewStorage(ctx context.Context, dsn string) (*Store, error) {

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Store{
		db:             pool,
		UserStore:      &UserStore{db: pool},
		PostStore:      &PostStore{db: pool},
		CommentStore:   &CommentStore{db: pool},
		FollowerStore:  &FollowerStore{db: pool},
		AuthTokenStore: &AuthTokenStore{db: pool},
		RoleStore:      &RoleStore{db: pool},
	}, nil

}
