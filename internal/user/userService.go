package user

import (
	"context"
	"fmt"
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
		return "", fmt.Errorf("CreateUser: error hashing password: %w", err)
	}

	user.Password = hashedPwd
	res, err := s.userStore.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("CreateUser: %w", err)
	}

	// Generate JWT
	jwtExp := time.Minute * 15
	accToken, _, err := s.jwtGenerator.GenerateToken(res.Id, res.Email, jwtExp)
	if err != nil {
		return "", fmt.Errorf("CreateUser: failed generating token: %w", err)
	}

	return accToken, nil
}

func (s *UserService) LoginUser(ctx context.Context, email string, password string) (*UserLoginRes, error) {
	var (
		accTokenExp     = time.Minute * 15
		refreshTokenExp = time.Hour * 24
	)

	user, err := s.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("LoginUser: %w", err)
	}

	ok, err := util.VerifyPassword(user.Password, password)
	if !ok {
		return nil, fmt.Errorf("LoginUser: failed verifying password: %w", err)
	}

	// Generate access token
	accToken, accClaim, err := s.jwtGenerator.GenerateToken(user.Id, user.Email, accTokenExp)
	if err != nil {
		return nil, fmt.Errorf("LoginUser: failed generating token: %w", err)
	}

	// Generate refresh token
	refreshToken, refreshClaim, err := s.jwtGenerator.GenerateToken(user.Id, user.Email, refreshTokenExp)

	userSession := Session{
		Id:           refreshClaim.RegisteredClaims.ID,
		UserEmail:    user.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaim.RegisteredClaims.ExpiresAt.Time,
	}

	session, err := s.userStore.CreateSession(ctx, &userSession)
	if err != nil {
		return nil, fmt.Errorf("LoginUser: failed creating session: %w", err)
	}

	res := &UserLoginRes{
		SessionId:       session.Id,
		AccessToken:     accToken,
		RefreshToken:    refreshToken,
		AccessTokenExp:  accClaim.ExpiresAt.Time,
		RefreshTokenExp: refreshClaim.ExpiresAt.Time,
		User: UserRes{
			Name:  user.Name,
			Email: user.Email,
		},
	}

	return res, nil
}

func (s *UserService) LogoutUser(ctx context.Context, id string) error {
	err := s.userStore.DeleteSession(ctx, id)
	if err != nil {
		fmt.Errorf("LogoutUser: failed deleting session: %w", err)
	}

	return nil
}

func (s *UserService) RenewAccessToken(ctx context.Context, refreshToken string) (*RenewAccessTokenRes, error) {
	refreshClaim, err := s.jwtGenerator.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("RenewAccessToken: failed verifying token: %w", err)
	}

	session, err := s.userStore.GetSession(ctx, refreshClaim.RegisteredClaims.ID)
	if err != nil {
		return nil, fmt.Errorf("RenewAccessToken: failed getting session: %w", err)
	}

	if session.IsRevoked {
		return nil, fmt.Errorf("RenewAccessToken: Session Revoked")
	}

	if session.UserEmail != refreshClaim.Email {
		return nil, fmt.Errorf("RenewAccessToken: Invalid session")
	}

	accToken, accClaims, err := s.jwtGenerator.GenerateToken(refreshClaim.Id, refreshClaim.Email, time.Minute*15)
	if err != nil {
		return nil, fmt.Errorf("RenewAccessToken: failed generating token: %w", err)
	}

	res := &RenewAccessTokenRes{
		AccessToken:          accToken,
		AccessTokenExpiresAt: accClaims.ExpiresAt.Time,
	}

	return res, nil
}
