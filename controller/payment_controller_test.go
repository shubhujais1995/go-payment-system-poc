package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"poc/model"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockPaymentMethodService is a mock implementation of the PaymentMethodService interface.
type MockPaymentMethodService struct {
	mock.Mock
}

func (m *MockPaymentMethodService) CreatePaymentMethod(db *gorm.DB, paymentMethod model.PaymentMethod) error {
	args := m.Called(db, paymentMethod)
	return args.Error(0)
}

// CreatePaymentMethodHandlerWithMock handles the creation of a new payment method with a mock service.
func CreatePaymentMethodHandlerWithMock(ctx iris.Context, service *MockPaymentMethodService) {
	var paymentMethodRequest PaymentMethodRequest

	// Parse request body
	if err := ctx.ReadJSON(&paymentMethodRequest); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]string{"error": "Invalid request body"})
		return
	}

	// Extract user ID
	payerID := ctx.Values().GetString("UserID")
	if payerID == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "User not authenticated"})
		return
	}

	// Map request to model
	paymentMethod := model.PaymentMethod{
		PayerID:       payerID,
		MethodType:    paymentMethodRequest.MethodType,
		CardNumber:    paymentMethodRequest.CardNumber,
		ExpiryDate:    paymentMethodRequest.ExpiryDate,
		Status:        paymentMethodRequest.Status,
		AccountNumber: paymentMethodRequest.AccountNumber,
		Details:       paymentMethodRequest.Details,
	}

	// Mock DB instance
	db := &gorm.DB{}

	// Call the service to create the payment method
	if err := service.CreatePaymentMethod(db, paymentMethod); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Success response
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(map[string]string{"message": "Payment method created successfully"})
}

func TestCreatePaymentMethodHandler(t *testing.T) {
	app := iris.New()

	tests := []struct {
		name               string
		requestBody        interface{}
		expectedStatusCode int
		expectedResponse   map[string]string
		userID             string
		mockServiceError   error
	}{
		{
			name: "Valid request",
			requestBody: PaymentMethodRequest{
				MethodType: "card",
				CardNumber: "1234567812345678",
				ExpiryDate: "12/23",
				Status:     "active",
				Details:    "Test card",
			},
			expectedStatusCode: iris.StatusCreated,
			expectedResponse:   map[string]string{"message": "Payment method created successfully"},
			userID:             "user123",
			mockServiceError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockPaymentMethodService)

			// Adjust the mock to expect *gorm.DB and model.PaymentMethod
			mockService.On("CreatePaymentMethod", mock.AnythingOfType("*gorm.DB"), mock.AnythingOfType("model.PaymentMethod")).Return(tt.mockServiceError)

			// Set up the handler with mock service
			handler := func(ctx iris.Context) {
				CreatePaymentMethodHandlerWithMock(ctx, mockService) // Pass mockService
			}
			app.Post("/payment-method", handler)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/payment-method", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rec := httptest.NewRecorder()

			// Create Iris context
			ctx := app.ContextPool.Acquire(rec, req)
			defer app.ContextPool.Release(ctx)

			// Set user ID
			if tt.userID != "" {
				ctx.Values().Set("UserID", tt.userID)
			}

			// Call handler
			handler(ctx)

			// Assert status code
			assert.Equal(t, tt.expectedStatusCode, rec.Code)

			// Assert response body
			var response map[string]string
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedResponse, response)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}
