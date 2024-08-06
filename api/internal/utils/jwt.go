package utils

import (
	"errors"
	"github.com/dunzane/brainbank-file/api/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserId int64
	Email  string
}

func GenerateToken(config config.Config, userId int64, email string) (string, error) {
	tokenExpiryHour := time.Duration(config.Auth.AccessExpire) * time.Minute
	tokenSecretKey := []byte(config.Auth.AccessSecret)

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiryHour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dunzane",
			Subject:   "brainbank",
			ID:        "1",
		},
		UserId: userId,
		Email:  email,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString(tokenSecretKey)
}

func ParseToken(config config.Config, tokenStr string) (*UserClaims, error) {
	tokenSecretKey := []byte(config.Auth.AccessSecret)

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return tokenSecretKey, nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("unknown claims type, cannot proceed")
	}
}
