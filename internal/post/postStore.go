package post

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostgrePostStore struct {
	db *sql.DB
}

func NewPostgrePostStorage(db *sql.DB) *PostgrePostStore {
	return &PostgrePostStore{
		db: db,
	}
}

func (s *PostgrePostStore) GetAllPost(ctx context.Context, userId int) (*[]PostRes, error) {
	var (
		post     PostRes
		postList []PostRes
	)

	query := "SELECT id, title, description, created_at FROM posts WHERE user_id = $1"
	res, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("[STORE][GetAllPost]: failed db query: %w", err)
	}

	defer res.Close()

	for res.Next() {
		if err := res.Scan(&post.Id, &post.Title, &post.Description, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("[STORE][GetAllPost]: failed scanning sql row: %w", err)
		}
		postList = append(postList, post)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("[STORE][GetAllPost]: failed iteration: %w", err)
	}

	return &postList, nil
}

func (s *PostgrePostStore) GetPostById(ctx context.Context, postId int) (*Post, error) {
	var post Post
	query := "SELECT id, title, description, created_at, updated_at FROM posts WHERE id=$1"
	err := s.db.QueryRowContext(ctx, query, postId).Scan(
		&post.Id,
		&post.Title,
		&post.Description,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("[STORE][GetPostById]: failed db query: %w", err)
	}

	return &post, err
}

func (s *PostgrePostStore) CreatePost(ctx context.Context, post *Post, userId int) error {
	query := "INSERT INTO posts (title, description, user_id) VALUES($1, $2, $3)"
	_, err := s.db.ExecContext(ctx, query, post.Title, post.Description, userId)
	if err != nil {
		return fmt.Errorf("[STORE][CreatePost]: failed db query: %w", err)
	}

	return nil
}

func (s *PostgrePostStore) EditPost(ctx context.Context, post *Post, userId int) (*Post, error) {
	query := "UPDATE posts SET title=$1, description=$2, updated_at=$3 WHERE id=$4 and user_id=$5 RETURNING id, title, description, created_at, updated_at, user_id"
	err := s.db.QueryRowContext(ctx, query, post.Title, post.Description, time.Now(), post.Id, userId).Scan(
		&post.Id,
		&post.Title,
		&post.Description,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.UserId,
	)

	if err != nil {
		return nil, fmt.Errorf("[STORE][EditPost]: failed db query: %w", err)
	}

	return post, nil
}
