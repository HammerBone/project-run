package user

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/HammerBone/project-run/internal/util"
)

type UserService struct {
	jwtGenerator util.JWTGenerator
	userStore    UserStore
}

func NewUserService(scrtKey string, userStore UserStore) *UserService {
	return &UserService{
		jwtGenerator: *util.NewJWTGenerator(scrtKey),
		userStore:    userStore,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *User) (string, error) {
	hashedPwd, err := util.HashPassword(user.Password)
	if err != nil {
		return "", err
	}

	user.Password = hashedPwd
	res, err := s.userStore.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	// Generate JWT
	jwtExp := time.Minute * 15
	accToken, _, err := s.jwtGenerator.GenerateToken(res.Id, res.Email, res.IsAdmin, jwtExp)
	if err != nil {
		return "", err
	}
	log.Printf("[CreateUser] Generated Access Token | accToken: %s", accToken)

	return accToken, nil
}

func (s *UserService) LoginUser(ctx context.Context, email string, password string) (*UserLoginRes, error) {
	var (
		accTokenExp     = time.Minute * 15
		refreshTokenExp = time.Hour * 24
		// userSessionExp  = time.Hour * 1
	)

	user, err := s.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	ok, err := util.VerifyPassword(user.Password, password)
	if !ok {
		return nil, err
	}

	// Generate access token
	accToken, accClaim, err := s.jwtGenerator.GenerateToken(user.Id, user.Email, user.IsAdmin, accTokenExp)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, refreshClaim, err := s.jwtGenerator.GenerateToken(user.Id, user.Email, user.IsAdmin, refreshTokenExp)

	userSession := Session{
		Id:           refreshClaim.RegisteredClaims.ID,
		UserEmail:    user.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaim.RegisteredClaims.ExpiresAt.Time,
	}

	session, err := s.userStore.CreateSession(ctx, &userSession)
	if err != nil {
		return nil, err
	}

	res := &UserLoginRes{
		SessionId:       session.Id,
		AccessToken:     accToken,
		RefreshToken:    refreshToken,
		AccessTokenExp:  accClaim.ExpiresAt.Time,
		RefreshTokenExp: refreshClaim.ExpiresAt.Time,
		User: UserRes{
			Name:    user.Name,
			Email:   user.Email,
			IsAdmin: user.IsAdmin,
		},
	}

	return res, nil
}

func (s *UserService) LogoutUser(ctx context.Context, id string) error {
	err := s.userStore.DeleteSession(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) RenewAccessToken(ctx context.Context, refreshToken string) (*RenewAccessTokenRes, error) {
	refreshClaim, err := s.jwtGenerator.VerifyToken(refreshToken)
	if err != nil {
		return nil, err
	}

	session, err := s.userStore.GetSession(ctx, refreshClaim.RegisteredClaims.ID)
	if err != nil {
		return nil, err
	}

	if session.IsRevoked {
		return nil, fmt.Errorf("Session Revoked")
	}

	if session.UserEmail != refreshClaim.Email {
		return nil, fmt.Errorf("Invalid session")
	}

	accToken, accClaims, err := s.jwtGenerator.GenerateToken(refreshClaim.Id, refreshClaim.Email, refreshClaim.IsAdmin, time.Minute*15)
	if err != nil {
		return nil, err
	}

	res := &RenewAccessTokenRes{
		AccessToken: accToken,
		AccessTokenExpiresAt: accClaims.ExpiresAt.Time,
	}

	return res, err
}