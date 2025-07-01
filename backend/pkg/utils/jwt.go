package utils

import (
	"fmt"
	"go-backend-v2/global"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for the given user ID
func GenerateToken(userID string) (string, error) {
	// Create claims with user ID and standard claims
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

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(global.Config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID
func ValidateToken(tokenString string) (string, error) {
	// Parse token with claims
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(global.Config.JWT.Secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	// Validate token and extract claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", fmt.Errorf("invalid token claims")
}

// ExtractUserIDFromToken extracts user ID from JWT token without full validation
// Used when we just need the user ID for logging or debugging
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

// IsTokenExpired checks if a token is expired without full validation
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
