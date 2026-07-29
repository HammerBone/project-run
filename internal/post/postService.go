package post

import "context"

type PostService struct {
	postStore *PostStore
}

func NewPostService(postStore *PostStore) *PostService {
	return &PostService{
		postStore: postStore,
	}
}

func (s *PostService) CreatePost(ctx context.Context, post *Post, userId int) error {
	err := s.postStore.CreatePost(ctx, post, userId)
	if err != nil {
		return err
	}

	return nil
}
