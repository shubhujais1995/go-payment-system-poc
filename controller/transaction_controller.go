package controller

import (
	"fmt"
	"poc/model"
	"poc/services"

	"github.com/kataras/iris/v12"
)

type TransactionHandler struct {
	svc *services.TransactionService
}

func CreateTransactionHandler(svc *services.TransactionService, ctx iris.Context) {
	payerId := ctx.Values().GetString("UserID")
	if payerId == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "User not authenticated"})
		return
	}

	// Parse request body
	var req model.ProcessPaymentInput
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]string{"error": "Invalid request payload"})
		return
	}

	fmt.Println("pay detail", req.PaymentDetails.CardNumber, "status", req.Status,
		"payment method id", req.PaymentDetails.ExpiryDate, "transaction id", req.PaymentDetails.CVV)

	// Create the transaction
	reservedAmount := 0.0
	transaction, err := svc.InitializeTransaction(ctx, payerId, req.PayeeID, req.Amount, req.TransactionType, req.Status, reservedAmount, req.PaymentMethodID, req.PaymentDetails)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	// ctx.JSON(transaction)
	ctx.JSON(iris.Map{
		"transaction_id": transaction.TransactionID,
		"status":         transaction.Status,
		"message":        "Transaction Completed Successfully.",
	})
}

// ListTransactionsHandler lists all transactions for the authenticated user
func ListTransactionsHandler(svc *services.TransactionService, ctx iris.Context) {
	// Extract authenticated user's ID
	userID := ctx.Values().GetString("UserID")
	if userID == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]string{"error": "User not authenticated"})
		return
	}

	// Fetch transactions
	// transactions, err := svc.ListTransactions(ctx.Request().Context(), userID)
	// if err != nil {
	// 	ctx.StatusCode(iris.StatusInternalServerError)
	// 	ctx.JSON(map[string]string{"error": err.Error()})
	// 	return
	// }

	ctx.JSON(nil)
}
