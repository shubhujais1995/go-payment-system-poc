package utils

import (
	"testing"
)

// Covered Cases:
// Valid Integer Strings

// "123", "456", "7890" → ✅ Should return true
// Non-Numeric Cases

// "12.34" → ❌ Decimal numbers
// "abc" → ❌ Alphabetic string
// "" → ❌ Empty string

// Missing Edge Cases
// Case	Reason	Expected
// "-123"	Negative integer	true
// "+123"	Positive sign prefix	true
// "000123"	Leading zeros	true
// " 123 "	Whitespace padding	false
// "123abc"	Numbers with letters	false
// "123 "	Trailing whitespace	false
// "1_000"	Underscore in numbers	false
// "123\n"	Newline character in numbers	false

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// ✅ Valid integer strings
		{"123", true},
		{"456", true},
		{"7890", true},

		// ❌ Invalid numeric formats
		{"12.34", false},  // Decimal number
		{"abc", false},    // Alphabetic characters
		{"", false},       // Empty string
		{"123abc", false}, // Mixed alphanumeric
		{"123 ", false},   // Trailing space
		{" 123", false},   // Leading space
		{"1_000", false},  // Underscore
		{"123\n", false},  // Newline character

		// ✅ Valid special cases
		{"-123", true},   // Negative number
		{"+123", true},   // Explicit positive number
		{"000123", true}, // Leading zeros
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := IsNumeric(test.input)
			if result != test.expected {
				t.Errorf("IsNumeric(%q) = %v; want %v", test.input, result, test.expected)
			}
		})
	}
}
