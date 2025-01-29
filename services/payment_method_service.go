package services

import (
	"errors"
	"poc/model"
	"poc/utils"
	"time"

	"gorm.io/gorm"
)

type PaymentMethodService struct {
	DB *gorm.DB
}

func NewPaymentMethodService(db *gorm.DB) *PaymentMethodService {
	return &PaymentMethodService{DB: db}
}

func (s *PaymentMethodService) CreatePaymentMethod(paymentMethod model.PaymentMethod) error {

	if paymentMethod.PaymentMethodID == "" {
		paymentMethod.PaymentMethodID = utils.GenerateUniqueID()
	}

	exists, err := s.CheckPaymentMethodExists(paymentMethod.PayerID, paymentMethod.MethodType, paymentMethod.CardNumber, paymentMethod.AccountNumber)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("payment method already exists for this payer")
	}

	switch paymentMethod.MethodType {
	case "card":
		if len(paymentMethod.CardNumber) != 16 || !utils.IsNumeric(paymentMethod.CardNumber) {
			return errors.New("invalid card number, must be 16 digits")
		}
		if paymentMethod.ExpiryDate == "" || len(paymentMethod.ExpiryDate) != 5 || !utils.IsValidExpiryDate(paymentMethod.ExpiryDate) {
			return errors.New("invalid expiry date, format should be MM/YY")
		}

	case "bank_transfer":
		if paymentMethod.AccountNumber == "" || len(paymentMethod.AccountNumber) < 11 || len(paymentMethod.AccountNumber) > 16 {
			return errors.New("invalid account number")
		}

	case "upi":
		if paymentMethod.Details == "" {
			return errors.New("UPI ID is required")
		}
		match, _ := utils.ValidateUpi(paymentMethod.Details)
		if !match {
			return errors.New("invalid UPI ID format, must be in the format example@upi")
		}

	case "wallet":
		if paymentMethod.Details == "" {
			return errors.New("wallet details are required")
		}
		if len(paymentMethod.Details) < 5 {
			return errors.New("wallet ID must be at least 5 characters long")
		}

	case "cheque":
		if paymentMethod.Details == "" {
			return errors.New("cheque number is required")
		}
		if !utils.IsNumeric(paymentMethod.Details) {
			return errors.New("cheque number must be numeric")
		}

	default:
		return errors.New("invalid payment method type")
	}

	paymentMethod.CreatedAt = time.Now()
	paymentMethod.UpdatedAt = time.Now()
	return s.DB.Create(&paymentMethod).Error
}

func (s *PaymentMethodService) CheckPaymentMethodExists(payerID, methodType, cardNumber, accountNumber string) (bool, error) {
	var existingPaymentMethod model.PaymentMethod
	var err error

	if methodType == "card" {
		err = s.DB.Where("payer_id = ? AND method_type = ? AND card_number = ?", payerID, methodType, cardNumber).First(&existingPaymentMethod).Error
		if err == nil {
			return true, nil
		}
	} else if methodType == "bank_transfer" {
		err = s.DB.Where("payer_id = ? AND method_type = ? AND account_number = ?", payerID, methodType, accountNumber).First(&existingPaymentMethod).Error
		if err == nil {
			return true, nil
		}
	} else {
		err = s.DB.Where("payer_id = ? AND method_type = ? AND (card_number = ? OR account_number = ?)", payerID, methodType, cardNumber, accountNumber).First(&existingPaymentMethod).Error
		if err == nil {
			return true, nil
		}
	}

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	return false, err
}

func (s *PaymentMethodService) GetPaymentMethods(payerID string) ([]model.PaymentMethod, error) {
	var paymentMethods []model.PaymentMethod
	err := s.DB.Where("payer_id = ?", payerID).Find(&paymentMethods).Error
	if err != nil {
		return nil, err
	}
	return paymentMethods, nil
}
func (s *PaymentMethodService) UpdatePaymentMethod(paymentMethodID string, updates map[string]interface{}) error {
	return s.DB.Model(&model.PaymentMethod{}).Where("payment_method_id = ?", paymentMethodID).Updates(updates).Error
}

func (s *PaymentMethodService) ValidatePaymentMethod(paymentMethodID string) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	err := s.DB.Where("payment_method_id = ?", paymentMethodID).First(&paymentMethod).Error
	if err != nil {
		return nil, err
	}
	return &paymentMethod, nil
}

func (s *PaymentMethodService) GetPaymentMethod(paymentMethodID string) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	err := s.DB.Where("payment_method_id = ?", paymentMethodID).First(&paymentMethod).Error
	if err != nil {
		return nil, err
	}
	return &paymentMethod, nil
}

func (s *PaymentMethodService) ValidatePaymentDetails(paymentMethod *model.PaymentMethod, paymentDetail model.PaymentDetails) error {
	if paymentMethod == nil {
		return errors.New("payment method is nil")
	}

	switch paymentMethod.MethodType {
	case "card":
		if paymentMethod.CardNumber != paymentDetail.CardNumber {
			return errors.New("payment method is not correct - card number")
		}
		if paymentMethod.ExpiryDate != paymentDetail.ExpiryDate {
			return errors.New("payment method is not correct - expiry date")
		}
	case "bank_transfer":
		if paymentMethod.AccountNumber != paymentDetail.AccountNumber {
			return errors.New("payment method is not correct - account number")
		}
	case "upi":
		if paymentMethod.Details != paymentDetail.UPIID {
			return errors.New("payment method is not correct - upi id")
		}
	case "wallet":
		if paymentMethod.Details != paymentDetail.Wallet {
			return errors.New("payment method is not correct - wallet")
		}
	case "cheque":
		if paymentMethod.Details != paymentDetail.Cheque {
			return errors.New("payment method is not correct - cheque")
		}
	default:
		return errors.New("invalid payment method type")
	}
	return nil
}
