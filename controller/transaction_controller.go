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

	// Create the transaction
	reservedAmount := 0.0
	transaction, err := svc.InitializeTransaction(ctx, payerId, req.PayeeID, req.Amount, req.TransactionType, req.Status, reservedAmount, req.TransactionID, req.PaymentMethodID, req.PaymentDetails)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(iris.Map{
		"transaction_id": transaction.TransactionID,
		"status":         transaction.Status,
		"message":        "Transaction Completed Successfully.",
	})
}

// type Transaction struct {
// 	TransactionID string
// 	PaymentMethod string
// 	Amount        float64
// 	Date          string
// }

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
	transactions, err := svc.ListTransactions(ctx.Request().Context(), userID)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]string{"error": err.Error()})
		return
	}

	// Define a slice to store the simplified results
	var transactionDetails []map[string]string

	// Loop through each transaction and extract only the desired fields
	for _, t := range transactions {
		transactionDetails = append(transactionDetails, map[string]string{
			"transaction_id":    t.TransactionID,
			"payment_method_id": t.PaymentMethodID,
			"transaction_type":  t.TransactionType,
			"status":            t.Status,
			"payer_name":        t.Payer.Name,
			"amount":            fmt.Sprintf("%.2f", t.Amount),
		})
	}

	ctx.JSON(transactionDetails)
}
