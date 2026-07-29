package post

import "context"

type PostRepository interface {
	CreatePost(ctx context.Context, post *Post) error
}