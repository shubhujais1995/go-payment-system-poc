package services

import (
	"fmt"
	"reflect"
	"testing"

	"poc/model"

	"gorm.io/gorm"
)

func TestGetPayee(t *testing.T) {
	mockDB := new(MockDB)
	payeeID := "test-payee-id"
	payee := model.Payee{
		PayeeID: payeeID,
		Name:    "Test Payee",
		Balance: 100.0,
	}

	mockDB.On("First", &model.Payee{}, "payee_id = ?", payeeID).Return(&payee, nil)
	mockDB.On("First", &model.Payee{}, "payee_id = ?", "non-existent-payee-id").Return(nil, gorm.ErrRecordNotFound)

	t.Run("Payee exists", func(t *testing.T) {
		mockDB := &gorm.DB{}
		result, err := getPayee(mockDB, payeeID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result.PayeeID != payeeID {
			t.Errorf("Expected payeeID %s, got %s", payeeID, result.PayeeID)
		}
	})

	t.Run("Payee does not exist", func(t *testing.T) {
		mockDB := &gorm.DB{}
		nonExistentPayeeID := "non-existent-payee-id"
		_, err := getPayee(mockDB, nonExistentPayeeID)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		expectedError := fmt.Sprintf("failed to retrieve payee with PayeeID %s: %v", nonExistentPayeeID, gorm.ErrRecordNotFound)
		if err.Error() != expectedError {
			t.Errorf("Expected error %v, got %v", expectedError, err)
		}
	})
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(append([]interface{}{dest}, conds...)...)
	if args.Get(0) != nil {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(args.Get(0)).Elem())
	}
	return &gorm.DB{Error: args.Error(1)}
}
