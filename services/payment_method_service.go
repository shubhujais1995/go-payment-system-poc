package services

import (
	"errors"
	"fmt"
	"poc/model"
	"poc/utils"
	"time"

	"gorm.io/gorm"
)

// PaymentMethodService provides methods for working with payment methods
type PaymentMethodService struct {
	DB *gorm.DB
}

// NewPaymentMethodService creates a new instance of PaymentMethodService
func NewPaymentMethodService(db *gorm.DB) *PaymentMethodService {
	return &PaymentMethodService{DB: db}
}

func (s *PaymentMethodService) CreatePaymentMethod(paymentMethod model.PaymentMethod) error {

	if paymentMethod.PaymentMethodID == "" {
		paymentMethod.PaymentMethodID = utils.GenerateUniqueID()
	}

	// Call CheckPaymentMethodExists to validate if payment method already exists
	exists, err := s.CheckPaymentMethodExists(paymentMethod.PayerID, paymentMethod.MethodType, paymentMethod.CardNumber, paymentMethod.AccountNumber)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("payment method already exists for this payer")
	}
	fmt.Println(paymentMethod, "-paymentMethod")
	// Validate MethodType
	if paymentMethod.MethodType == "" {
		return errors.New("payment method type is required")
	}

	// Validate based on MethodType
	switch paymentMethod.MethodType {
	case "card":
		// Validate CardNumber and ExpiryDate for card payment method
		if len(paymentMethod.CardNumber) != 16 || !utils.IsNumeric(paymentMethod.CardNumber) {
			return errors.New("invalid card number, must be 16 digits")
		}
		if paymentMethod.ExpiryDate == "" || len(paymentMethod.ExpiryDate) != 5 || !utils.IsValidExpiryDate(paymentMethod.ExpiryDate) {
			return errors.New("invalid expiry date, format should be MM/YY")
		}

	case "bank_transfer":
		// Validate AccountNumber for bank transfer
		if paymentMethod.AccountNumber == "" || len(paymentMethod.AccountNumber) < 11 || len(paymentMethod.AccountNumber) > 16 {
			return errors.New("invalid account number")
		}

	case "upi":
		// Validate UPI ID format
		if paymentMethod.Details == "" {
			return errors.New("UPI ID is required")
		}
		// UPI ID format validation (example@upi)
		match, _ := utils.ValidateUpi(paymentMethod.Details)
		if !match {
			return errors.New("invalid UPI ID format, must be in the format example@upi")
		}

	case "wallet":
		// Validate Wallet details (it could be an ID, number, etc.)
		if paymentMethod.Details == "" {
			return errors.New("wallet details are required")
		}
		// Optional: You could add length checks or format checks for wallet IDs, if needed
		if len(paymentMethod.Details) < 5 { // Example validation: Wallet ID should be at least 5 characters
			return errors.New("wallet ID must be at least 5 characters long")
		}

	case "cheque":
		// Validate Cheque number
		if paymentMethod.Details == "" {
			return errors.New("cheque number is required")
		}
		// Optional: Cheque number validation, for example checking if it's numeric
		if !utils.IsNumeric(paymentMethod.Details) {
			return errors.New("cheque number must be numeric")
		}

	default:
		return errors.New("invalid payment method type")
	}

	// Insert payment method into the database
	paymentMethod.CreatedAt = time.Now()
	paymentMethod.UpdatedAt = time.Now()
	return s.DB.Create(&paymentMethod).Error
}

// CheckPaymentMethodExists checks if a payment method already exists for the given payer.
func (s *PaymentMethodService) CheckPaymentMethodExists(payerID, methodType, cardNumber, accountNumber string) (bool, error) {
	var existingPaymentMethod model.PaymentMethod
	var err error
	// Adjust the query based on the method_type
	if methodType == "card" {
		// Only check for card_number
		err := s.DB.Where("payer_id = ? AND method_type = ? AND card_number = ?", payerID, methodType, cardNumber).First(&existingPaymentMethod).Error
		if err == nil {
			// Payment method already exists for this card
			return true, nil
		}
	} else if methodType == "bank_transfer" {
		// Only check for account_number
		err := s.DB.Where("payer_id = ? AND method_type = ? AND account_number = ?", payerID, methodType, accountNumber).First(&existingPaymentMethod).Error
		if err == nil {
			// Payment method already exists for this account number
			return true, nil
		}
	} else {
		// For other types, check both fields (if required)
		err := s.DB.Where("payer_id = ? AND method_type = ? AND (card_number = ? OR account_number = ?)", payerID, methodType, cardNumber, accountNumber).First(&existingPaymentMethod).Error
		if err == nil {
			// Payment method already exists
			return true, nil
		}
	}

	// If the record is not found, it's valid to add a new one
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	// Handle other DB errors
	return false, err
}

/* one
// CreatePaymentMethod validates and creates a payment method
func (s *PaymentMethodService) CreatePaymentMethod(paymentMethod model.PaymentMethod) error {

	// Call CheckPaymentMethodExists to validate if payment method already exists
	exists, err := s.CheckPaymentMethodExists(paymentMethod.PayerID, paymentMethod.MethodType, paymentMethod.CardNumber, paymentMethod.AccountNumber)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("payment method already exists for this payer")
	}

	// Validate MethodType
	// fmt.Fprintln(paymentMethod.MethodType, " = paymentMethod")
	if paymentMethod.MethodType == "" {
		return errors.New("payment method type is required")
	}

	// Validate based on MethodType
	switch paymentMethod.MethodType {
	case "card":
		// Validate CardNumber and ExpiryDate for card payment method
		if len(paymentMethod.CardNumber) != 16 || !utils.IsNumeric(paymentMethod.CardNumber) {
			return errors.New("invalid card number, must be 16 digits")
		}
		if paymentMethod.ExpiryDate == "" || len(paymentMethod.ExpiryDate) != 5 || !utils.IsValidExpiryDate(paymentMethod.ExpiryDate) {
			return errors.New("invalid expiry date, format should be MM/YY")
		}

	case "bank_transfer":
		// Validate AccountNumber for bank transfer
		if paymentMethod.AccountNumber == "" || len(paymentMethod.AccountNumber) < 8 {
			return errors.New("invalid account number")
		}

	case "upi":
		// Validate UPI if required
		if paymentMethod.Details == "" {
			return errors.New("UPI ID is required")
		}

	case "wallet":
		// Validate Wallet method if required
		if paymentMethod.Details == "" {
			return errors.New("wallet details are required")
		}

	case "cheque":
		// Validate for cheque method if required
		if paymentMethod.Details == "" {
			return errors.New("cheque number is required")
		}

	default:
		return errors.New("invalid payment method type")
	}

	// Insert payment method into the database
	paymentMethod.CreatedAt = time.Now()
	paymentMethod.UpdatedAt = time.Now()
	return s.DB.Create(&paymentMethod).Error
}

// CheckPaymentMethodExists checks if a payment method already exists for the given payer.
func (s *PaymentMethodService) CheckPaymentMethodExists(payerID, methodType, cardNumber, accountNumber string) (bool, error) {
	var existingPaymentMethod model.PaymentMethod

	// Adjust the query based on the method_type
	if methodType == "card" {
		// Only check for card_number
		err := s.DB.Where("payer_id = ? AND method_type = ? AND card_number = ?", payerID, methodType, cardNumber).First(&existingPaymentMethod).Error
		if err == nil {
			// Payment method already exists for this card
			return true, nil
		}
	} else if methodType == "bank_transfer" {
		// Only check for account_number
		err := s.DB.Where("payer_id = ? AND method_type = ? AND account_number = ?", payerID, methodType, accountNumber).First(&existingPaymentMethod).Error
		if err == nil {
			// Payment method already exists for this account number
			return true, nil
		}
	} else {
		// For other types, check both fields (if required)
		err := s.DB.Where("payer_id = ? AND method_type = ? AND (card_number = ? OR account_number = ?)", payerID, methodType, cardNumber, accountNumber).First(&existingPaymentMethod).Error
		if err == nil {
			// Payment method already exists
			return true, nil
		}
	}

	// If the record is not found, it's valid to add a new one
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	// Handle other DB errors
	return false, err
}
old working with already exist error
*/
// GetPaymentMethods fetches all payment methods for a given payer ID
func (s *PaymentMethodService) GetPaymentMethods(payerID string) ([]model.PaymentMethod, error) {
	var paymentMethods []model.PaymentMethod
	err := s.DB.Where("payer_id = ?", payerID).Find(&paymentMethods).Error
	if err != nil {
		return nil, err
	}
	return paymentMethods, nil
}
func (svc *TransactionService) VerifyFetchedPaymentMethod(paymentMethod *model.PaymentMethod) error {

	if paymentMethod == nil || paymentMethod.Status != "active" {
		return errors.New("invalid or inactive payment method")
	}
	return nil
}

// ==================================== ***** ============================= //
// import (
// 	"errors"
// 	"fmt"
// 	"poc/model"
// 	"time"

// 	"gorm.io/gorm"
// )

// type PaymentMethodService struct {
// 	DB *gorm.DB
// }

// func NewPaymentMethodService(db *gorm.DB) *PaymentMethodService {
// 	return &PaymentMethodService{DB: db}
// }

// func (s *PaymentMethodService) CreatePaymentMethod(paymentMethod model.PaymentMethod) error {
// 	// Ensure the payment method details are properly set
// 	if paymentMethod.PayerID == "" || paymentMethod.MethodType == "" || paymentMethod.Details == "" {
// 		return errors.New("payment method details are incomplete")
// 	}

// 	// Set creation and update timestamps
// 	paymentMethod.CreatedAt = time.Now()
// 	paymentMethod.UpdatedAt = time.Now()

// 	// Insert payment method into the database
// 	if err := s.DB.Create(&paymentMethod).Error; err != nil {
// 		return fmt.Errorf("failed to create payment method: %w", err)
// 	}

// 	return nil
// }

// // GetPaymentMethod fetches a payment method by payer ID
// func (s *PaymentMethodService) GetPaymentMethod(payerID string) (model.PaymentMethod, error) {
// 	var paymentMethod model.PaymentMethod
// 	err := s.DB.First(&paymentMethod, "payer_id = ?", payerID).Error
// 	if err != nil {
// 		return paymentMethod, err
// 	}
// 	return paymentMethod, nil
// }

// ********** NEED TO WORK ON IT LATER ******************//
// UpdatePaymentMethod updates the payment method details
func (s *PaymentMethodService) UpdatePaymentMethod(paymentMethodID string, updates map[string]interface{}) error {
	// Update the payment method in the database using the given updates
	updates["updated_at"] = time.Now()
	return s.DB.Model(&model.PaymentMethod{}).
		Where("payment_method_id = ?", paymentMethodID).
		Updates(updates).Error
}

// ValidatePaymentMethod ensures the payment method is valid and active
func (s *PaymentMethodService) ValidatePaymentMethod(paymentMethodID string) (model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	err := s.DB.First(&paymentMethod, "payment_method_id = ?", paymentMethodID).Error
	if err != nil {
		return paymentMethod, err
	}
	if paymentMethod.Status != "active" {
		return paymentMethod, errors.New("payment method is not active")
	}
	return paymentMethod, nil

	// // Fetch default payment method if not provided
	// if paymentMethodID == "" {
	// 	var paymentMethod model.PaymentMethod
	// 	err := s.DB.Where("payer_id = ?", payerID).First(&paymentMethod).Error
	// 	if err != nil {
	// 		return nil, errors.New("default payment method not found for user")
	// 	}
	// 	paymentMethodID = paymentMethod.PaymentMethodID
	// }

}

// package services

// import (
// 	"errors"
// 	"fmt"
// 	"poc/model"
// 	"poc/utils"
// 	"time"

// 	"gorm.io/gorm"
// )

// type PaymentMethodService struct {
// 	DB *gorm.DB
// }

// func NewPaymentMethodService(db *gorm.DB) *PaymentMethodService {
// 	return &PaymentMethodService{DB: db}
// }

// // CreatePaymentMethod adds a new payment method to the database
// func (s *PaymentMethodService) CreatePaymentMethod(payerID, methodType, details, expiryDate, cvv, status string) (*model.PaymentMethod, error) {
// 	// Generate a unique PaymentMethodID
// 	paymentMethodID := utils.GenerateUniqueID()

// 	pm := &model.PaymentMethod{
// 		PaymentMethodID: paymentMethodID,
// 		PayerID:         payerID,
// 		MethodType:      methodType,
// 		Details:         details,
// 		ExpiryDate:      expiryDate,
// 		CVV:             cvv,
// 		Status:          status,
// 		CreatedAt:       time.Now(),
// 		UpdatedAt:       time.Now(),
// 	}

// 	// Save the payment method to the database
// 	if err := s.DB.Create(pm).Error; err != nil {
// 		return nil, fmt.Errorf("failed to create payment method: %v", err)
// 	}

// 	return pm, nil
// }

// // GetPaymentMethodByID fetches a payment method by ID
// func (s *PaymentMethodService) GetPaymentMethodByID(paymentMethodID string) (model.PaymentMethod, error) {
// 	var pm model.PaymentMethod
// 	err := s.DB.First(&pm, "payment_method_id = ?", paymentMethodID).Error
// 	if err != nil {
// 		return pm, err
// 	}
// 	return pm, nil
// }

// // UpdatePaymentMethod updates the details or status of a payment method
// func (s *PaymentMethodService) UpdatePaymentMethod(paymentMethodID string, updates map[string]interface{}) error {
// 	return s.DB.Model(&model.PaymentMethod{}).Where("payment_method_id = ?", paymentMethodID).Updates(updates).Error
// }

// // ValidatePaymentMethod ensures the payment method is valid and active
// func (s *PaymentMethodService) ValidatePaymentMethod(paymentMethodID string) (model.PaymentMethod, error) {
// 	pm, err := s.GetPaymentMethodByID(paymentMethodID)
// 	if err != nil {
// 		return pm, err
// 	}
// 	if pm.Status != "active" {
// 		return pm, errors.New("payment method is not active")
// 	}
// 	return pm, nil
// }
