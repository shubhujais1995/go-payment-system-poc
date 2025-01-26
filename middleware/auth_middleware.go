package middleware

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kataras/iris/v12"
)

// Secret key for signing (should be the same as used during token generation)
var secretKey = []byte("your-secret-key")

// AuthMiddleware checks for the validity of JWT token
func AuthMiddleware(ctx iris.Context) {
	// Extract token from the "Authorization" header (Bearer token format)
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "Authorization header missing"})
		return
	}

	// Expect token in the form of "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "Invalid Authorization format"})
		return
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "Invalid or expired token"})
		return
	}

	// Check if token is valid
	if !token.Valid {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "Invalid token"})
		return
	}

	// Extract user ID from claims and store it in the context
	claims, ok := token.Claims.(jwt.MapClaims)
	// fmt.Println("claims - ", claims)
	if !ok || claims["UserID"] == nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "Invalid token claims"})
		return
	}

	// Set the user ID into the context for further use
	userID := claims["UserID"].(string)
	// fmt.Println("Authenticated user ID:", userID) // Log user ID for debugging
	ctx.Values().Set("UserID", userID)

	// Call the next handler
	ctx.Next()
}
