package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain-text password.
// Returns the hashed password or an error.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err) // Keep this for internal debugging
		return "", err                                // Return error to the caller for further handling
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	// Compare the provided password with the hash stored in the database
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		// If bcrypt returns an error, it means the password doesn't match
		return false
	}
	return true
}
