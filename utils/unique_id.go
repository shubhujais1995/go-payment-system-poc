package utils

import "github.com/google/uuid"

// GenerateUniqueID generates a new UUID string
func GenerateUniqueID() string {
	return uuid.NewString()
}
