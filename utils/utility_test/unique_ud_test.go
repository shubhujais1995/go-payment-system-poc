package utils

import (
	"poc/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateUniqueID(t *testing.T) {
	id1 := utils.GenerateUniqueID()
	id2 := utils.GenerateUniqueID()

	// Ensure the generated ID is a valid UUID
	_, err1 := uuid.Parse(id1)
	_, err2 := uuid.Parse(id2)

	assert.NoError(t, err1, "Generated ID1 is not a valid UUID")
	assert.NoError(t, err2, "Generated ID2 is not a valid UUID")

	// Ensure the IDs are unique
	assert.NotEqual(t, id1, id2, "Generated IDs should be unique")
}
