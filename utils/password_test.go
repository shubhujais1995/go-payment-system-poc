package utils

import "testing"

// Scenarios Covered:
// Test Case	Function Tested	Expected Outcome
// Hashing a valid password	HashPassword	Returns a non-empty hash
// Hashing an empty password	HashPassword	Returns a hash (bcrypt allows empty passwords)
// Hashing a long password	HashPassword	Returns a valid hash
// Comparing correct password with hash	CheckPasswordHash	Returns true
// Comparing incorrect password with hash	CheckPasswordHash	Returns false
// Comparing empty password with hash	CheckPasswordHash	Returns false
// Comparing hash with an empty string	CheckPasswordHash	Returns false
// Comparing hash with a completely invalid hash	CheckPasswordHash	Returns false

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid Password", "securepassword", false},
		{"Empty Password", "", false},
		{"Max Length Password (72 chars)", string(make([]byte, 72)), false}, // ✅ bcrypt's limit
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if hash == "" {
				t.Errorf("Expected a non-empty hash for %q", tt.name)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "securepassword"
	hash, _ := HashPassword(password) // Hash a known password

	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{"Correct Password", password, hash, true},
		{"Incorrect Password", "wrongpassword", hash, false},
		{"Empty Password", "", hash, false},
		{"Empty Hash", password, "", false},
		{"Completely Invalid Hash", password, "invalidhash", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPasswordHash(tt.password, tt.hash)
			if result != tt.expected {
				t.Errorf("CheckPasswordHash(%q, %q) = %v; want %v", tt.password, tt.hash, result, tt.expected)
			}
		})
	}
}
