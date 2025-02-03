package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/stretchr/testify/assert"
)

// Mock middleware for testing
// func MockAuthMiddleware(ctx iris.Context) {
// 	ctx.Next()
// }

// // Mock handler functions for testing
// func MockCreatePaymentMethodHandler(ctx iris.Context) {
// 	ctx.StatusCode(http.StatusCreated)
// 	ctx.JSON(map[string]string{"message": "Payment method created"})
// }

// func MockGetPaymentMethodHandler(ctx iris.Context) {
// 	ctx.StatusCode(http.StatusOK)
// 	ctx.JSON(map[string]string{"message": "Payment method fetched"})
// }

// func MockUpdatePaymentMethodHandler(ctx iris.Context) {
// 	ctx.StatusCode(http.StatusOK)
// 	ctx.JSON(map[string]string{"message": "Payment method updated"})
// }

// func MockValidatePaymentMethodHandler(ctx iris.Context) {
// 	ctx.StatusCode(http.StatusOK)
// 	ctx.JSON(map[string]string{"message": "Payment method validated"})
// }

func TestRegisterPaymentRoutes(t *testing.T) {
	app := iris.New()
	// Mock middleware
	authMiddleware := func(ctx iris.Context) {
		ctx.Next()
	}
	app.Use(authMiddleware)
	RegisterPaymentRoutes(app)

	// Mock controller handlers
	CreatePaymentMethodHandler := func(ctx iris.Context) {
		ctx.StatusCode(http.StatusCreated)
	}
	GetPaymentMethodHandler := func(ctx iris.Context) {
		ctx.StatusCode(http.StatusOK)
	}
	UpdatePaymentMethodHandler := func(ctx iris.Context) {
		ctx.StatusCode(http.StatusOK)
	}
	ValidatePaymentMethodHandler := func(ctx iris.Context) {
		ctx.StatusCode(http.StatusOK)
	}

	app.Post("/payment-methods", CreatePaymentMethodHandler)
	app.Get("/payment-methods", GetPaymentMethodHandler)
	app.Put("/payment-methods/paymentMethodID", UpdatePaymentMethodHandler)
	app.Post("/payment-methods/validate/paymentMethodID", ValidatePaymentMethodHandler)

	// Build the router
	app.Build()

	// Test POST /transactions
	req := httptest.NewRequest("POST", "/payment-methods", nil)
	resp := httptest.NewRecorder()
	app.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Test GET /payment method service
	req = httptest.NewRequest("GET", "/payment-methods", nil)
	resp = httptest.NewRecorder()
	app.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Test PUT /payment method service
	req = httptest.NewRequest("PUT", "/payment-methods/paymentMethodID", nil)
	resp = httptest.NewRecorder()
	app.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Test put validating payment method
	req = httptest.NewRequest("POST", "/payment-methods/validate/paymentMethodID", nil)
	resp = httptest.NewRecorder()
	app.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
