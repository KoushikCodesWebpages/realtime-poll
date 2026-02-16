package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var wsSecret = []byte("WS_SUPER_SECRET_123") // separate secret

func GenerateWSToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(5 * time.Minute).Unix(), // short lived
		"iat": time.Now().Unix(),
		"type": "ws",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(wsSecret)
}

func ParseWSToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid signing method")
		}
		return wsSecret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	if claims["type"] != "ws" {
		return "", errors.New("not ws token")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid subject")
	}

	return userID, nil
}