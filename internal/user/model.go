package user

import "time"

type User struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserRes struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserSignUpRes struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	AccToken string `json:"access_token"`
}

type UserLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRes struct {
	SessionId       string    `json:"session_id"`
	AccessToken     string    `json:"access_token"`
	RefreshToken    string    `json:"refresh_token"`
	AccessTokenExp  time.Time `json:"access_token_exp"`
	RefreshTokenExp time.Time `json:"refresh_token_exp"`
	User            UserRes   `json:"user"`
}

type UserLogoutReq struct {
	SessionId string `json:"session_id"`
}

type Session struct {
	Id           string    `json:"id"`
	UserEmail    string    `json:"user_email"`
	RefreshToken string    `json:"refresh_token"`
	IsRevoked    bool      `json:"is_revoked"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RenewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
