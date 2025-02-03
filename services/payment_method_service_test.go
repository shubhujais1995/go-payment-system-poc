package services

import (
	"poc/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// func TestValidatePaymentDetails(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		paymentMethod *model.PaymentMethod
// 		paymentDetail model.PaymentDetails
// 		expectedError error
// 	}{
// 		{
// 			name:          "Nil Payment Method",
// 			paymentMethod: nil,
// 			paymentDetail: model.PaymentDetails{},
// 			expectedError: errors.New("payment method is nil"),
// 		},
// 		{
// 			name: "Card Number Mismatch",
// 			paymentMethod: &model.PaymentMethod{
// 				MethodType: "card",
// 				CardNumber: "1234567812345678",
// 				ExpiryDate: "12/23",
// 			},
// 			paymentDetail: model.PaymentDetails{
// 				CardNumber: "8765432187654321",
// 				ExpiryDate: "12/23",
// 			},
// 			expectedError: errors.New("payment method is not correct - card number"),
// 		},
// 		{
// 			name: "Card Expiry Date Mismatch",
// 			paymentMethod: &model.PaymentMethod{
// 				MethodType: "card",
// 				CardNumber: "1234567812345678",
// 				ExpiryDate: "12/23",
// 			},
// 			paymentDetail: model.PaymentDetails{
// 				CardNumber: "1234567812345678",
// 				ExpiryDate: "11/23",
// 			},
// 			expectedError: errors.New("payment method is not correct - expiry date"),
// 		},
// 		{
// 			name: "Bank Account Number Mismatch",
// 			paymentMethod: &model.PaymentMethod{
// 				MethodType:    "bank_transfer",
// 				AccountNumber: "12345678901",
// 			},
// 			paymentDetail: model.PaymentDetails{
// 				AccountNumber: "10987654321",
// 			},
// 			expectedError: errors.New("payment method is not correct - account number"),
// 		},
// 		{
// 			name: "UPI ID Mismatch",
// 			paymentMethod: &model.PaymentMethod{
// 				MethodType: "upi",
// 				Details:    "example@upi",
// 			},
// 			paymentDetail: model.PaymentDetails{
// 				UPIID: "wrong@upi",
// 			},
// 			expectedError: errors.New("payment method is not correct - upi id"),
// 		},
// 		{
// 			name: "Wallet ID Mismatch",
// 			paymentMethod: &model.PaymentMethod{
// 				MethodType: "wallet",
// 				Details:    "wallet123",
// 			},
// 			paymentDetail: model.PaymentDetails{
// 				Wallet: "wallet321",
// 			},
// 			expectedError: errors.New("payment method is not correct - wallet"),
// 		},
// 		{
// 			name: "Cheque Number Mismatch",
// 			paymentMethod: &model.PaymentMethod{
// 				MethodType: "cheque",
// 				Details:    "123456",
// 			},
// 			paymentDetail: model.PaymentDetails{
// 				Cheque: "654321",
// 			},
// 			expectedError: errors.New("payment method is not correct - cheque"),
// 		},
// 		{
// 			name: "Invalid Payment Method Type",
// 			paymentMethod: &model.PaymentMethod{
// 				MethodType: "invalid",
// 			},
// 			paymentDetail: model.PaymentDetails{},
// 			expectedError: errors.New("invalid payment method type"),
// 		},
// 		{
// 			name: "Valid Card Payment Method",
// 			paymentMethod: &model.PaymentMethod{
// 				MethodType: "card",
// 				CardNumber: "1234567812345678",
// 				ExpiryDate: "12/23",
// 			},
// 			paymentDetail: model.PaymentDetails{
// 				CardNumber: "1234567812345678",
// 				ExpiryDate: "12/23",
// 			},
// 			expectedError: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := ValidatePaymentDetails(tt.paymentMethod, tt.paymentDetail)
// 			if tt.expectedError != nil {
// 				assert.EqualError(t, err, tt.expectedError.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 		})
// 	}
// }

// // MockDB is a mock implementation of the gorm.DB.
// type MockDB struct {
// 	mock.Mock
// }

// func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
// 	args = m.Called(query, args)
// 	return args.Get(0).(*gorm.DB)
// }

// func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
// 	args := m.Called(out, where)
// 	return args.Get(0).(*gorm.DB)
// }

// func TestGetPaymentMethod(t *testing.T) {
// 	tests := []struct {
// 		name                  string
// 		paymentMethodID       string
// 		mockPaymentMethod     *model.PaymentMethod
// 		mockError             error
// 		expectedPaymentMethod *model.PaymentMethod
// 		expectedError         error
// 	}{
// 		{
// 			name:            "Payment Method Found",
// 			paymentMethodID: "valid_id",
// 			mockPaymentMethod: &model.PaymentMethod{
// 				PaymentMethodID: "valid_id",
// 				MethodType:      "card",
// 				CardNumber:      "1234567812345678",
// 				ExpiryDate:      "12/23",
// 			},
// 			mockError: nil,
// 			expectedPaymentMethod: &model.PaymentMethod{
// 				PaymentMethodID: "valid_id",
// 				MethodType:      "card",
// 				CardNumber:      "1234567812345678",
// 				ExpiryDate:      "12/23",
// 			},
// 			expectedError: nil,
// 		},
// 		{
// 			name:                  "Payment Method Not Found",
// 			paymentMethodID:       "invalid_id",
// 			mockPaymentMethod:     nil,
// 			mockError:             gorm.ErrRecordNotFound,
// 			expectedPaymentMethod: nil,
// 			expectedError:         gorm.ErrRecordNotFound,
// 		},
// 		{
// 			name:                  "Database Error",
// 			paymentMethodID:       "error_id",
// 			mockPaymentMethod:     nil,
// 			mockError:             errors.New("database error"),
// 			expectedPaymentMethod: nil,
// 			expectedError:         errors.New("database error"),
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Mock the database call
// 			db := new(MockDB)
// 			db.On("Where", "payment_method_id = ?", tt.paymentMethodID).Return(db)
// 			db.On("First", &model.PaymentMethod{}).Return(tt.mockPaymentMethod, tt.mockError)

// 			initializer.SetDB(db)

// 			paymentMethod, err := GetPaymentMethod(tt.paymentMethodID)
// 			if tt.expectedError != nil {
// 				assert.EqualError(t, err, tt.expectedError.Error())
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 			assert.Equal(t, tt.expectedPaymentMethod, paymentMethod)
// 		})
// 	}
// }

// MockDB is a mock implementation of the gorm.DB.
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	callArgs := m.Called(out, where)
	return callArgs.Get(0).(*gorm.DB)
}

func TestGetPaymentMethods(t *testing.T) {
	mockDB := new(MockDB)
	payerID := "payer123"
	expectedPaymentMethods := []model.PaymentMethod{
		{
			PaymentMethodID: "1",
			PayerID:         payerID,
			MethodType:      "card",
			CardNumber:      "1234567812345678",
			ExpiryDate:      "12/23",
		},
		{
			PaymentMethodID: "2",
			PayerID:         payerID,
			MethodType:      "bank_transfer",
			AccountNumber:   "12345678901",
		},
	}

	// Mock the DB responses
	mockDB.On("Where", "payer_id = ?", []interface{}{payerID}).Return(mockDB)
	mockDB.On("Find", &[]model.PaymentMethod{}, []interface{}{}).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*[]model.PaymentMethod)
		*arg = expectedPaymentMethods
	}).Return(mockDB)

	// Mock the initializer.GetDB function to return the mockDB
	// originalGetDB := initializer.GetDB()
	// defer func() { initializer.SetDB(originalGetDB) }()
	// initializer.SetDB(mockDB)

	// Call the GetPaymentMethods function
	paymentMethods, err := GetPaymentMethods(payerID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, paymentMethods)
	assert.Equal(t, len(expectedPaymentMethods), len(paymentMethods))
	assert.Equal(t, expectedPaymentMethods[0].PaymentMethodID, paymentMethods[0].PaymentMethodID)
	assert.Equal(t, expectedPaymentMethods[1].PaymentMethodID, paymentMethods[1].PaymentMethodID)

	// Verify mock expectations
	mockDB.AssertExpectations(t)
}

// func TestGetPaymentMethods_NotFound(t *testing.T) {
// 	mockDB := new(MockDB)
// 	payerID := "payer123"

// 	// Mock the DB responses
// 	mockDB.On("Where", "payer_id = ?", []interface{}{payerID}).Return(mockDB)
// 	mockDB.On("Find", &[]model.PaymentMethod{}, []interface{}{}).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})

// 	// Mock the initializer.GetDB function to return the mockDB
// 	// initializer.GetDB = func() *gorm.DB {
// 	// 	return mockDB
// 	// }

// 	// Call the GetPaymentMethods function
// 	paymentMethods, err := GetPaymentMethods(payerID)

// 	// Assertions
// 	assert.NoError(t, err)
// 	assert.Empty(t, paymentMethods)

// 	// Verify mock expectations
// 	mockDB.AssertExpectations(t)
// }

// func TestGetPaymentMethods_DBError(t *testing.T) {
// 	mockDB := new(MockDB)
// 	payerID := "payer123"
// 	expectedError := errors.New("database error")

// 	// Mock the DB responses
// 	mockDB.On("Where", "payer_id = ?", []interface{}{payerID}).Return(mockDB)
// 	mockDB.On("Find", &[]model.PaymentMethod{}, []interface{}{}).Return(&gorm.DB{Error: expectedError})

// 	// Mock the initializer.GetDB function to return the mockDB
// 	// initializer.GetDB = func() *gorm.DB {
// 	// 	return mockDB
// 	// }

// 	// Call the GetPaymentMethods function
// 	paymentMethods, err := GetPaymentMethods(payerID)

// 	// Assertions
// 	assert.Error(t, err)
// 	assert.Nil(t, paymentMethods)
// 	assert.Equal(t, expectedError, err)

// 	// Verify mock expectations
// 	mockDB.AssertExpectations(t)
// }
