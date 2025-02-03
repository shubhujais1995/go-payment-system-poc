package controller

import (
	"fmt"
	"log"
	"poc/initializer"
	"poc/model"
	"poc/services"

	"github.com/kataras/iris/v12"
)

type PaymentMethodService interface {
	CreatePaymentMethod(paymentMethod model.PaymentMethod) error
}

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
func CreatePaymentMethodHandler(ctx iris.Context) {
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
	db := initializer.GetDB() // Get DB instance once
	// Call the service to create the payment method
	if err := services.CreatePaymentMethod(db, paymentMethod); err != nil {
		log.Printf("Error creating payment method: %v", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Respond with success
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(map[string]string{"message": "Payment method created successfully"})
}
func GetPaymentMethodHandler(ctx iris.Context) {
	payerID := ctx.Values().GetString("UserID")

	if payerID == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "User not authenticated"})
		return
	}

	paymentMethods, err := services.GetPaymentMethods(payerID)
	if err != nil {
		log.Printf("Error fetching payment methods: %v", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": "Could not fetch payment methods"})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(paymentMethods)
}

func UpdatePaymentMethodHandler(ctx iris.Context) {
	paymentMethodID := ctx.Params().GetString("paymentMethodID")
	var updates map[string]interface{}

	if err := ctx.ReadJSON(&updates); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request body"})
		return
	}

	if err := services.UpdatePaymentMethod(paymentMethodID, updates); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Payment method updated successfully"})
}

func ValidatePaymentMethodHandler(ctx iris.Context) {
	paymentMethodID := ctx.Params().GetString("paymentMethodID")

	paymentMethod, err := services.ValidatePaymentMethod(paymentMethodID)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Payment method is valid", "payment_method": paymentMethod})
}
