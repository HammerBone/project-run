package post

import "context"

type PostStore interface {
	GetAllPost(ctx context.Context, userId int) (*[]PostRes, error)
	GetPostById(ctx context.Context, postId int) (*Post, error)
	CreatePost(ctx context.Context, post *Post, userId int) error
	EditPost(ctx context.Context, post *Post, userId int) (*Post, error)
}
