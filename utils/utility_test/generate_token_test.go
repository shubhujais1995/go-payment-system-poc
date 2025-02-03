package utils

import (
	"fmt"
	"poc/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestGenerateToken(t *testing.T) {
	userID := "12345"
	token, err := utils.GenerateToken(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return utils.SecretKey, nil
	})

	if err != nil {
		t.Fatalf("Expected no error while parsing token, got %v", err)
	}

	if !parsedToken.Valid {
		t.Fatalf("Expected token to be valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("Expected claims to be of type jwt.MapClaims")
	}

	if claims["UserID"] != userID {
		t.Errorf("Expected UserID to be %v, got %v", userID, claims["UserID"])
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatalf("Expected exp to be a float64")
	}

	if time.Unix(int64(exp), 0).Sub(time.Now()) > time.Hour*72 {
		t.Errorf("Expected token to expire in 72 hours or less")
	}
}
