package utils

import (
	"poc/utils"
	"testing"
)

func TestIsValidExpiryDate(t *testing.T) {
	tests := []struct {
		expiryDate string
		expected   bool
	}{
		{"01/23", true},
		{"12/99", true},
		{"00/23", false},
		{"13/23", false},
		{"1/23", false},
		{"01/2023", false},
		{"01-23", false},
		{"", false},
	}

	for _, test := range tests {
		result := utils.IsValidExpiryDate(test.expiryDate)
		if result != test.expected {
			t.Errorf("IsValidExpiryDate(%s) = %v; expected %v", test.expiryDate, result, test.expected)
		}
	}
}
