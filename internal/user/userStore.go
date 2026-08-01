package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostgreUserStore struct {
	db *sql.DB
}

func NewPostgreUserStorage(db *sql.DB) *PostgreUserStore {
	return &PostgreUserStore{db: db}
}

func (s *PostgreUserStore) CreateUser(ctx context.Context, user *User) (*User, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return nil, fmt.Errorf("CreateUser: failed loading time location: %w", err)
	}

	var (
		id        int64
		timeInWIB = time.Now().In(loc)
	)

	query := "INSERT INTO users (name, email, password, is_admin, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = s.db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password, user.IsAdmin, timeInWIB).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("CreateUser: failed db query: %w", err)
	}

	user.Id = int(id)

	return user, err
}

func (s *PostgreUserStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := "SELECT id, name, email, password, is_admin, created_at FROM users WHERE email=$1"
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("GetUserByEmail: failed db query: %w", err)
	}

	return &user, nil
}

func (s *PostgreUserStore) CreateSession(ctx context.Context, session *Session) (*Session, error) {
	query := "INSERT INTO sessions (id, user_email, refresh_token, is_revoked, created_at, expires_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := s.db.ExecContext(
		ctx, query,
		session.Id,
		session.UserEmail,
		session.RefreshToken,
		session.IsRevoked,
		session.CreatedAt,
		session.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateSession: failed db query: %w", err)
	}

	return session, nil
}

func (s *PostgreUserStore) GetSession(ctx context.Context, id string) (*Session, error) {
	var ses Session

	query := "SELECT id, user_email, refresh_token, is_revoked, expires_at FROM sessions WHERE id=$1"
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&ses.Id,
		&ses.UserEmail,
		&ses.RefreshToken,
		&ses.IsRevoked,
		&ses.ExpiresAt,
	)

	if err != nil {
		return nil, fmt.Errorf("GetSession: failed db query: %w", err)
	}

	return &ses, nil
}

func (s *PostgreUserStore) DeleteSession(ctx context.Context, id string) error {
	query := "DELETE FROM sessions WHERE id=$1"
	_, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("DeleteSession: failed db query: %w", err)
	}

	return nil
}
