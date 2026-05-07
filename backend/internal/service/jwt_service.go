package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("dev-secret-key")

func GenerateJWT(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ParseJWT(tokenString string) (int, string, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", jwt.ErrTokenInvalidClaims
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", jwt.ErrTokenInvalidClaims
	}

	email, ok := claims["email"].(string)
	if !ok {
		return 0, "", jwt.ErrTokenInvalidClaims
	}

	return int(userIDFloat), email, nil
}