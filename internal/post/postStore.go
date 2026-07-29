package post

import (
	"context"
	"database/sql"
)

type PostStore struct {
	db *sql.DB
}

func NewPostgrePostStorage(db *sql.DB) *PostStore {
	return &PostStore{
		db: db,
	}
}

func (s *PostStore) CreatePost(ctx context.Context, post *Post, userId int) error {
	query := "INSERT INTO posts (description, user_id) VALUES($1, $2)"
	_, err := s.db.ExecContext(ctx, query, post.Description, userId)
	if err != nil {
		return err
	}

	return nil
}
