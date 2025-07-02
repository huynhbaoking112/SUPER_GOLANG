package utils

import (
	"fmt"
	"go-backend-v2/global"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(global.Config.JWT.ExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-backend-v2",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(global.Config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(global.Config.JWT.Secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", fmt.Errorf("invalid token claims")
}

func ExtractUserIDFromToken(tokenString string) string {
	token, _ := jwt.ParseWithClaims(tokenString, &JWTClaims{}, nil)
	if token == nil {
		return ""
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims.UserID
	}

	return ""
}

func IsTokenExpired(tokenString string) bool {
	token, _ := jwt.ParseWithClaims(tokenString, &JWTClaims{}, nil)
	if token == nil {
		return true
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims.ExpiresAt.Before(time.Now())
	}

	return true
}
