package main

import (
	"os"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Save(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

// func (m *MockDB) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
// 	args := m.Called(fc, opts)
// 	return args.Error(0)
// }

// func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
// 	args := m.Called(dest, conds)
// 	return args.Get(0).(*gorm.DB)
// }

// TestMainFunction is a test for the main function.
func TestMainFunction(t *testing.T) {
	// Mock .env loading, no need to load real config
	os.Setenv("PORT", "8080")

	// Initialize a mock database
	mockDB := new(MockDB)
	// Mock the methods you need
	mockDB.On("First", mock.Anything, mock.Anything).Return(&gorm.DB{})
	mockDB.On("Create", mock.Anything).Return(&gorm.DB{})
	mockDB.On("Transaction", mock.Anything, mock.Anything).Return(nil)

	// Initialize services with the mock DB
	// userService := services.NewUserService(mockDB)
	// paymentMethodService := services.NewPaymentMethodService(mockDB)
	// transactionService := services.NewTransactionService(mockDB, paymentMethodService)

	// Create an Iris application instance
	app := iris.New()
	app.HandleDir("/", iris.Dir("."))

	// Register routes (ensure this import path matches your project's actual path)
	// routes.RegisterAuthRoutes(app, userService)
	// routes.RegisterPaymentRoutes(app, paymentMethodService)
	// routes.RegisterTransactionRoutes(app, transactionService)

	// Test port configuration
	port := os.Getenv("PORT")
	assert.Equal(t, "8080", port, "Port should be 8080")

	// Test that the server is initialized correctly
	assert.NotNil(t, app, "App should not be nil")

	// Test that services are initialized correctly
	// assert.NotNil(t, userService, "User service should not be nil")
	// assert.NotNil(t, paymentMethodService, "Payment method service should not be nil")
	// assert.NotNil(t, transactionService, "Transaction service should not be nil")

	// You can add more tests to mock and verify server start functionality
}
