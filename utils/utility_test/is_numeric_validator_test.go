package utils

import (
	"poc/utils"
	"testing"
)

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"456", true},
		{"7890", true},
		{"12.34", false},
		{"abc", false},
		{"", false},
	}

	for _, test := range tests {
		result := utils.IsNumeric(test.input)
		if result != test.expected {
			t.Errorf("IsNumeric(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
