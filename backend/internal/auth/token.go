package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManagerInterface interface {
	IssueAccessToken(user User) (AccessToken, error)
}

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

type AccessToken struct {
	Value     string
	ExpiresAt time.Time
}

func NewTokenManager(secret []byte, ttl time.Duration) *TokenManager {
	return &TokenManager{secret: secret, ttl: ttl}
}

func (tm *TokenManager) IssueAccessToken(user User) (AccessToken, error) {
	now := time.Now()
	expiresAt := now.Add(tm.ttl)

	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.Id, 10),
		Issuer:    "authlearn",
		Audience:  jwt.ClaimStrings{"authlearn-api"},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(tm.secret)

	if err != nil {
		return AccessToken{}, err
	}

	return AccessToken{signedToken, expiresAt}, nil
}
