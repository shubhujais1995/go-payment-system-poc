package utils

import (
	"testing"
)

func TestIsValidExpiryDate(t *testing.T) {
	tests := []struct {
		expiryDate string
		expected   bool
	}{
		// ✅ Valid cases
		{"01/23", true},
		{"12/99", true},

		// ❌ Invalid month
		{"00/23", false},
		{"13/23", false},

		// ❌ Invalid format
		{"1/23", false},
		{"01/2023", false},
		{"01-23", false},
		{"", false},

		// ❗ Additional edge cases
		{"00/00", false},   // Both month & year invalid
		{"11/aa", false},   // Non-numeric year
		{"ab/23", false},   // Non-numeric month
		{"01/2", false},    // Year has only one digit
		{" 01/23 ", false}, // Spaces before/after
		{"01//23", false},  // Extra `/` in format
	}

	for _, test := range tests {
		t.Run(test.expiryDate, func(t *testing.T) {
			result := IsValidExpiryDate(test.expiryDate)
			if result != test.expected {
				t.Errorf("IsValidExpiryDate(%s) = %v; expected %v", test.expiryDate, result, test.expected)
			}
		})
	}
}

func TestValidateUpi(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		// ✅ Positive Cases
		{"Valid UPI ID", "user@bank", true, false},
		{"Valid UPI with dot", "test.user@upi", true, false},
		{"Valid UPI with numbers", "test123@upi", true, false},
		{"Valid UPI with hyphen", "test-user@upi", true, false},
		{"Valid UPI with underscore", "test_user@upi", true, false},

		// ❌ Negative Cases
		{"Empty String", "", false, false},
		{"Missing @ symbol", "userbank", false, false},
		{"Invalid Characters", "user!@bank", false, false},
		{"Only Domain", "@bank", false, false},
		{"Space in UPI", "user @bank", false, false},
		{"Multiple @ symbols", "user@bank@upi", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateUpi(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpi() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateUpi() = %v, want %v", got, tt.want)
			}
		})
	}
}
