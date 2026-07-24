package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type userSecretClaim struct {
	Id      int
	Email   string
	IsAdmin bool
	jwt.RegisteredClaims
}

type JWTGenerator struct {
	scrtKey string
}

func NewJWTGenerator(scrtKey string) *JWTGenerator {
	return &JWTGenerator{
		scrtKey: scrtKey,
	}
}

func (j *JWTGenerator) NewUserClaim(id int, email string, isAdmin bool, duration time.Duration) (*userSecretClaim, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("New user claim error error: %+v", err)
	}

	return &userSecretClaim{
		Id:      id,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}

func (j *JWTGenerator) GenerateToken(id int, email string, isAdmin bool, duration time.Duration) (string, *userSecretClaim, error) {
	userClaim, err := j.NewUserClaim(id, email, isAdmin, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	tokenStr, err := token.SignedString([]byte(j.scrtKey))
	if err != nil {
		return "", nil, err
	}

	return tokenStr, userClaim, nil
}

func (j *JWTGenerator) VerifyToken(tokenStr string) (*userSecretClaim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, userSecretClaim{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("[VerifyToken] Invalid signing method")
		}
		return j.scrtKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("[VerifyToken] error parsing token: %+v", err)
	}

	claim, ok := token.Claims.(*userSecretClaim)
	if !ok {
		return nil, fmt.Errorf("[VerifyToken] Invalid claim")
	}

	return claim, err
}
