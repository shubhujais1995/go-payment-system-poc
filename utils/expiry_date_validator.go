package utils

import "regexp"

// Helper function to validate expiry date in MM/YY format
func IsValidExpiryDate(expiryDate string) bool {
	re := regexp.MustCompile(`^(0[1-9]|1[0-2])\/([0-9]{2})$`)
	return re.MatchString(expiryDate)
}

func ValidateUpi(paymentMethodDetails string) (bool, error) {
	upiRegex := `^[a-zA-Z0-9_.-]+@[a-zA-Z0-9.-]+$`

	return regexp.MatchString(upiRegex, paymentMethodDetails)
}
