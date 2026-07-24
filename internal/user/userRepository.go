package user

import (
	"context"
)

type UserStore interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateSession(ctx context.Context, session *Session) (*Session, error)
	GetSession(ctx context.Context, id string) (*Session, error)
	DeleteSession(ctx context.Context, id string) error
}
