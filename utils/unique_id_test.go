
package utils

import (
	"regexp"
	"testing"
)
// ✅ Ensuring the function returns a non-empty string
// ✅ Checking the format of the UUID using regex
// ✅ Verifying that multiple calls produce unique values

// Regular expression to validate a UUID (version 4)
var uuidRegex = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)

func TestGenerateUniqueID(t *testing.T) {
	// Generate two UUIDs
	id1 := GenerateUniqueID()
	id2 := GenerateUniqueID()

	// ✅ Check if the generated ID is not empty
	if id1 == "" {
		t.Errorf("Expected non-empty UUID, got empty string")
	}

	// ✅ Validate the UUID format using regex
	if !uuidRegex.MatchString(id1) {
		t.Errorf("Generated UUID %q does not match expected format", id1)
	}

	// ✅ Ensure that multiple calls generate unique values
	if id1 == id2 {
		t.Errorf("Expected unique UUIDs, but got duplicates: %q and %q", id1, id2)
	}
}
