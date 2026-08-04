package post

import "context"

type PostStore interface {
	GetPostById(ctx context.Context, postId int) (*Post, error)
	CreatePost(ctx context.Context, post *Post, userId int) error
	EditPost(ctx context.Context, post *Post, userId int) (*Post, error)
}
