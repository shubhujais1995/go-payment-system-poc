package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key for signing the token (use a secure key in production)
var secretKey = []byte("your-secret-key")

// GenerateToken generates a JWT token
func GenerateToken(userID string) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"UserID": userID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("Failed to generate token: %v", err)
	}

	return signedToken, nil
}
