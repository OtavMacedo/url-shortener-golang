package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	secretKey  []byte
	expiration time.Duration
}

func NewAuthService(secretKey string, expiration time.Duration) *AuthService {
	return &AuthService{
		secretKey:  []byte(secretKey),
		expiration: expiration,
	}
}

func (as *AuthService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(as.expiration).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(as.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate jwt: %w", err)
	}
	return signedToken, nil
}
