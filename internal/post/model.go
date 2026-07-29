package post

import "time"

type Post struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	UserId      int    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
