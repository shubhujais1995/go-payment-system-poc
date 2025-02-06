package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Positive Cases (Already Covered)
// Valid Token Generation
// Ensures GenerateToken returns a token without an error.
// Validates the token signature and claims.

// Negative Cases (Missing & Should Be Added)
// Case	Reason	Expected Outcome
// Empty userID	User ID is required for token generation	Should return a valid token but empty user ID in claims
// Token Tampering	If the token is modified, it should be invalid	Parsing should fail with an error
// Invalid Secret Key	If a wrong key is used to verify the token	Should return an error
// Expired Token	If the token is manually modified to an expired timestamp	Should return an error
// Unsupported Algorithm	If a different algorithm is used	Parsing should fail

// TestGenerateToken tests the token generation function
func TestGenerateToken(t *testing.T) {
	// ✅ Positive Case: Generate valid token
	userID := "12345"
	token, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Parse the token to verify it
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return SecretKey, nil
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

// 🔴 Negative Test Cases

// 1️⃣ Test with Empty UserID
func TestGenerateToken_EmptyUserID(t *testing.T) {
	token, err := GenerateToken("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if token == "" {
		t.Errorf("Expected a valid token, got an empty string")
	}
}

// 2️⃣ Test with Token Tampering
func TestParseTamperedToken(t *testing.T) {
	userID := "12345"
	token, _ := GenerateToken(userID)

	// Modify the token (tampering)
	tamperedToken := token + "invalid"

	_, err := jwt.Parse(tamperedToken, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err == nil {
		t.Errorf("Expected error for tampered token, got nil")
	}
}

// 3️⃣ Test with Invalid Secret Key
func TestParseWithInvalidSecretKey(t *testing.T) {
	userID := "12345"
	token, _ := GenerateToken(userID)

	// Attempt to parse with a different secret key
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("wrong-secret-key"), nil
	})

	if err == nil {
		t.Errorf("Expected error for invalid secret key, got nil")
	}
}

// 4️⃣ Test with Expired Token
func TestParseExpiredToken(t *testing.T) {
	// Generate a token with an expired timestamp
	expiredClaims := jwt.MapClaims{
		"UserID": "12345",
		"exp":    time.Now().Add(-time.Hour).Unix(), // Token already expired
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	signedToken, _ := token.SignedString(SecretKey)

	// Attempt to parse
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err == nil || parsedToken.Valid {
		t.Errorf("Expected expired token to be invalid, got valid token")
	}
}

// 5️⃣ Test with Unsupported Algorithm
func TestParseUnsupportedAlgorithm(t *testing.T) {
	// Generate token with a different signing method (RS256 instead of HS256)
	token := jwt.New(jwt.SigningMethodRS256)
	signedToken, _ := token.SignedString(SecretKey)

	_, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err == nil {
		t.Errorf("Expected error for unsupported algorithm, got nil")
	}
}
