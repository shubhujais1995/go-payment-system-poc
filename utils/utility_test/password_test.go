package utils

import (
	"poc/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	t.Run("Valid password hashing", func(t *testing.T) {
		password := "securepassword123"
		hashedPassword, err := utils.HashPassword(password)
		assert.NoError(t, err, "Expected no error while hashing password")
		assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")

		// Verify the hashed password is valid
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		assert.NoError(t, err, "Expected the hashed password to match the original password")
	})

	t.Run("Empty password", func(t *testing.T) {
		hashedPassword, err := utils.HashPassword("")
		assert.NoError(t, err, "Hashing an empty password should not return an error")
		assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
	})
}

func TestCheckPasswordHash(t *testing.T) {
	password := "SecurePass123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err, "Error generating hash")

	t.Run("Valid Password", func(t *testing.T) {
		match := utils.CheckPasswordHash(password, string(hashedPassword))
		assert.True(t, match, "Password should match")
	})

	t.Run("Invalid Password", func(t *testing.T) {
		match := utils.CheckPasswordHash("WrongPass", string(hashedPassword))
		assert.False(t, match, "Password should not match")
	})
}
