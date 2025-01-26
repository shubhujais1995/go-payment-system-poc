package utils

import "strconv"

// Helper function to check if a string is numeric
func IsNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
