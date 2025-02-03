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

// MockTransactionService is a mock implementation of the TransactionService interface.
type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) InitializeTransaction(db *gorm.DB, ctx iris.Context, payerId, payeeId string, amount float64, transactionType, status string, reservedAmount float64, transactionId, paymentMethodId string, paymentDetails model.PaymentDetails) (*model.Transaction, error) {
	args := m.Called(ctx, payerId, payeeId, amount, transactionType, status, reservedAmount, transactionId, paymentMethodId, paymentDetails)
	return args.Get(0).(*model.Transaction), args.Error(1)
}

// CreateTransactionHandlerWithMock handles the creation of a transaction using a mock service.
func CreateTransactionHandlerWithMock(ctx iris.Context, service *MockTransactionService) {
	var transactionRequest model.ProcessPaymentInput

	// Parse request body
	if err := ctx.ReadJSON(&transactionRequest); err != nil {
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

	// Mock DB instance
	db := &gorm.DB{}

	// Call the service to initialize the transaction
	transaction, err := service.InitializeTransaction(db, ctx, payerID, transactionRequest.PayeeID, transactionRequest.Amount, transactionRequest.TransactionType, transactionRequest.Status, 0.0, "", "", transactionRequest.PaymentDetails)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Success response
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(map[string]string{"transaction_id": transaction.TransactionID, "status": transaction.Status, "message": "Transaction Completed Successfully."})
}

func TestCreateTransactionHandler(t *testing.T) {
	app := iris.New()

	tests := []struct {
		name               string
		requestBody        interface{}
		expectedStatusCode int
		expectedResponse   map[string]string
		userID             string
		mockServiceError   error
		mockTransaction    *model.Transaction
	}{
		{
			name: "Valid request",
			requestBody: model.ProcessPaymentInput{
				PayeeID:         "payee123",
				Amount:          100.0,
				TransactionType: "payment",
				Status:          "completed",
				PaymentDetails:  model.PaymentDetails{CardNumber: "1234567812345678"},
			},
			expectedStatusCode: iris.StatusCreated,
			expectedResponse:   map[string]string{"transaction_id": "txn123", "status": "completed", "message": "Transaction Completed Successfully."},
			userID:             "user123",
			mockServiceError:   nil,
			mockTransaction: &model.Transaction{
				TransactionID: "txn123",
				Status:        "completed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockTransactionService)

			// Adjust the mock to expect the correct parameters
			mockService.On("InitializeTransaction", mock.Anything, tt.userID, tt.requestBody.(model.ProcessPaymentInput).PayeeID, tt.requestBody.(model.ProcessPaymentInput).Amount, tt.requestBody.(model.ProcessPaymentInput).TransactionType, tt.requestBody.(model.ProcessPaymentInput).Status, 0.0, "", "", tt.requestBody.(model.ProcessPaymentInput).PaymentDetails).Return(tt.mockTransaction, tt.mockServiceError)

			// Set up the handler with mock service
			handler := func(ctx iris.Context) {
				CreateTransactionHandlerWithMock(ctx, mockService)
			}
			app.Post("/transaction", handler)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewReader(body))
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
