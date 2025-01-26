package controller

import (
	"fmt"
	"log"
	"poc/model"
	"poc/services"

	"github.com/kataras/iris/v12"
)

// PaymentMethodRequest represents the request payload for creating a payment method
type PaymentMethodRequest struct {
	MethodType    string `json:"method_type" validate:"required,oneof=card bank_transfer upi wallet cheque"`
	CardNumber    string `json:"card_number,omitempty" validate:"required_if=MethodType card,len=16"`
	ExpiryDate    string `json:"expiry_date,omitempty" validate:"required_if=MethodType card"`
	Status        string `json:"status" validate:"required,oneof=active inactive"`
	Details       string `json:"details" validate:"required"`
	AccountNumber string `json:"account_number,omitempty" validate:"required_if=MethodType bank_transfer"`
}

// CreatePaymentMethodHandler handles the creation of a new payment method
func CreatePaymentMethodHandler(svc *services.PaymentMethodService, ctx iris.Context) {
	var paymentMethodRequest PaymentMethodRequest

	// Parse the request body into the paymentMethodRequest struct
	if err := ctx.ReadJSON(&paymentMethodRequest); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]string{"error": "Invalid request body"})
		return
	}

	// Extract user ID from JWT token (already done by AuthMiddleware)
	payerID := ctx.Values().GetString("UserID") // This will be set by your AuthMiddleware

	// If no userID is found, return an error (middleware should ensure this)
	if payerID == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "User not authenticated"})
		return
	}

	// Map PaymentMethodRequest to PaymentMethod model
	paymentMethod := model.PaymentMethod{
		PayerID:       payerID,
		MethodType:    paymentMethodRequest.MethodType,
		CardNumber:    paymentMethodRequest.CardNumber,
		ExpiryDate:    paymentMethodRequest.ExpiryDate,
		Status:        paymentMethodRequest.Status,
		AccountNumber: paymentMethodRequest.AccountNumber,
		Details:       paymentMethodRequest.Details,
	}

	fmt.Println(paymentMethod, " = paymentMethod")

	// Call the service to create the payment method
	if err := svc.CreatePaymentMethod(paymentMethod); err != nil {
		log.Printf("Error creating payment method: %v", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Respond with success
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(map[string]string{"message": "Payment method created successfully"})
}

// GetPaymentMethodsHandler handles fetching payment methods for a specific payer
func GetPaymentMethodHandler(svc *services.PaymentMethodService, ctx iris.Context) {
	// Extract user ID from JWT token (already done by AuthMiddleware)
	payerID := ctx.Values().GetString("UserID") // This will be set by your AuthMiddleware

	// If no userID is found, return an error (middleware should ensure this)
	if payerID == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "User not authenticated"})
		return
	}

	// Call the service to fetch payment methods
	paymentMethods, err := svc.GetPaymentMethods(payerID)
	if err != nil {
		log.Printf("Error fetching payment methods: %v", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": "Could not fetch payment methods"})
		return
	}

	// Respond with the payment methods
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(paymentMethods)
}

/*
// PaymentMethodRequest represents the incoming request body for creating a payment method.
type PaymentMethodRequest struct {
	MethodType string `json:"method_type" validate:"required"`
	Details    string `json:"details" validate:"required"`
	ExpiryDate string `json:"expiry_date,omitempty"`
	Status     string `json:"status" validate:"required"`
	PayerID    string `json:"payer_id,omitempty"` // We'll set this manually
}

func CreatePaymentMethodHandler(svc *services.PaymentMethodService, ctx iris.Context) {
	// Extract user ID from JWT token (already done by AuthMiddleware)
	userID := ctx.Values().GetString("UserID") // This will be set by your AuthMiddleware

	// If no userID is found, return an error (middleware should ensure this)
	if userID == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "User not authenticated"})
		return
	}

	// Parse payment method data from the request body
	var paymentMethodData PaymentMethodRequest
	if err := ctx.ReadJSON(&paymentMethodData); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]string{"error": "Invalid request body"})
		return
	}

	// Manually set the PayerID from the authenticated user
	paymentMethodData.PayerID = userID
	// Map the PaymentMethodRequest data to the model.PaymentMethod struct
	paymentMethod := model.PaymentMethod{
		PayerID:    paymentMethodData.PayerID,
		MethodType: paymentMethodData.MethodType,
		Details:    paymentMethodData.Details,
		ExpiryDate: paymentMethodData.ExpiryDate,
		Status:     paymentMethodData.Status,
	}

	fmt.Println(paymentMethod, "Payment Method Data")

	// Call the service to create the payment method
	if err := svc.CreatePaymentMethod(paymentMethod); err != nil {
		log.Printf("Error creating payment method: %v", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": "Could not create payment method"})
		return
	}

	// Respond with success
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(map[string]string{"message": "Payment method created successfully"})
}

// GetPaymentMethodHandler handles fetching the payment method for the authenticated user
func GetPaymentMethodHandler(svc *services.PaymentMethodService, ctx iris.Context) {
	// Retrieve the user ID from the token (authentication middleware sets this)
	payerID := ctx.Values().GetString("UserID")

	// Get the payment method for the user from the service
	paymentMethod, err := svc.GetPaymentMethod(payerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"error": "Payment method not found"})
		} else {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Unable to fetch payment method"})
		}
		return
	}

	// Respond with the payment method details
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(paymentMethod)
}
*/

// ============= ********** ====================//
// Need to work on it later

// UpdatePaymentMethodHandler handles updating a payment method
func UpdatePaymentMethodHandler(svc *services.PaymentMethodService, ctx iris.Context) {
	paymentMethodID := ctx.Params().GetString("paymentMethodID")
	var updates map[string]interface{}

	// Bind the JSON request to the updates map
	if err := ctx.ReadJSON(&updates); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request body"})
		return
	}

	// Call the service to update the payment method
	if err := svc.UpdatePaymentMethod(paymentMethodID, updates); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Payment method updated successfully"})
}

// ValidatePaymentMethodHandler handles validating a payment method
func ValidatePaymentMethodHandler(svc *services.PaymentMethodService, ctx iris.Context) {
	paymentMethodID := ctx.Params().GetString("paymentMethodID")

	// Call the service to validate the payment method
	paymentMethod, err := svc.ValidatePaymentMethod(paymentMethodID)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Payment method is valid", "payment_method": paymentMethod})
}
