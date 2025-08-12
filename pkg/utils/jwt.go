package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(userID, username, role string, secret []byte) (string, error) {
	expTime := time.Now().Add(24 * time.Hour)

	mapClaims := jwt.MapClaims{
		"exp":      jwt.NewNumericDate(expTime),
		"iat":      jwt.NewNumericDate(time.Now()),
		"sub":      userID,
		"username": username,
		"role":     role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
