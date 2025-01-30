package services

import (
	"testing"

	"poc/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mocking gorm.DB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

// Test for getPayer method
func TestGetPayer_Success(t *testing.T) {
	// Create a mock DB
	mockDB := new(MockDB)

	// Define expected Payer data
	expectedPayer := model.Payer{
		PayerID: "payer123",
		Balance: 100.00,
	}

	// Mock the DB's First method to return the expected payer data
	mockDB.On("First", mock.Anything, "PayerID = ?", "payer123").Return(mockDB).Run(func(args mock.Arguments) {
		// Assign the result to the out argument (expectedPayer)
		*out.(*model.Payer) = expectedPayer
	})

	// Create the service with the mocked DB
	svc := &TransactionService{
		DB: mockDB,
	}

	// Call getPayer
	payer, err := svc.getPayer(mockDB, "payer123")

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, payer)
	assert.Equal(t, "payer123", payer.PayerID)
	assert.Equal(t, 100.00, payer.Balance)

	// Verify that the mock was called
	mockDB.AssertExpectations(t)
}
