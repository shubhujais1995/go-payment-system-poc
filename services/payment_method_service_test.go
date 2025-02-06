package services

import (
	"poc/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock implementation of the gorm.DB.
type MockDB struct {
	mock.Mock
	*gorm.DB
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
	t.Skip()
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
