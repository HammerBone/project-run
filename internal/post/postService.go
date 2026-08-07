package post

import (
	"context"
	"fmt"
)

type PostService struct {
	postStore PostStore
}

func NewPostService(postStore PostStore) *PostService {
	return &PostService{
		postStore: postStore,
	}
}

func (s *PostService) GetAllPost(ctx context.Context, userId int) (*ListPostRes, error) {
	res, err := s.postStore.GetAllPost(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE][GetAllPost]: %w", err)
	}

	return &ListPostRes{
		Post: *res,
	}, nil
}

func (s *PostService) CreatePost(ctx context.Context, post *Post, userId int) error {
	err := s.postStore.CreatePost(ctx, post, userId)
	if err != nil {
		return fmt.Errorf("[SERVICE][CreatePost]: %w", err)
	}

	return nil
}

func (s *PostService) EditPost(ctx context.Context, p *Post, userId int) (*Post, error) {
	res, err := s.postStore.EditPost(ctx, p, userId)
	if err != nil {
		return nil, fmt.Errorf("[SERVICE][EditPost]: %w", err)
	}

	return res, nil
}
