package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"poc/model"
	"poc/utils"

	"gorm.io/gorm"
)

var errorLogger *log.Logger

func init() {
	// Ensure the errors directory exists
	if err := os.MkdirAll("errors", os.ModePerm); err != nil {
		fmt.Println("Error creating errors directory:", err)
		return
	}

	// Open the log file inside the errors directory
	errorLogFile, err := os.OpenFile("errors/transaction_service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}

	// Initialize the error logger
	errorLogger = log.New(errorLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

type TransactionService struct {
	DB                   *gorm.DB
	PaymentMethodService *PaymentMethodService
}

func NewTransactionService(db *gorm.DB, pmService *PaymentMethodService) *TransactionService {
	return &TransactionService{
		DB:                   db,
		PaymentMethodService: pmService,
	}
}

func (svc *TransactionService) InitializeTransaction(ctx context.Context, payerID, payeeID string, amount float64, transactionType, status string, reservedAmount float64, transactionID, paymentMethodID string, paymentDetail model.PaymentDetails) (*model.Transaction, error) {
	var transaction *model.Transaction

	err := svc.DB.Transaction(func(tx *gorm.DB) error {
		// Step 1: Retrieve payer and payee
		payer, payee, err := svc.getPayerAndPayee(tx, payerID, payeeID)
		if err != nil {
			return err
		}

		// Step 2: Validate payment method
		if err := svc.validatePaymentMethod(ctx, payerID, paymentMethodID, paymentDetail); err != nil {
			return err
		}

		if err := svc.validateTransactionPayload(tx, transactionID, payerID, payeeID, amount, transactionType, paymentMethodID); err != nil {
			return err
		}

		// // Step 3: Check for duplicate transactions
		// if svc.isDuplicateTransaction(tx, payerID, payeeID, amount, transactionType) {
		// 	return errors.New("duplicate transaction detected")
		// }

		// Step 4: Check balance
		if !svc.hasSufficientBalance(payer, amount) {
			return errors.New("insufficient balance")
		}

		// Step 5: Create the transaction
		transaction, err = svc.createTransaction(tx, payer, payee, amount, transactionType, status, reservedAmount, paymentMethodID, paymentDetail)
		if err != nil {
			return err
		}

		// Step 6: Process payment
		if err := svc.processPayment(tx, transaction); err != nil {
			return err
		}

		// Step 7: Create audit entry
		if err := svc.createAuditEntry(tx, transaction); err != nil {
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return transaction, err
}

func (svc *TransactionService) getPayerAndPayee(tx *gorm.DB, payerID, payeeID string) (*model.Payer, *model.Payee, error) {
	payer, err := svc.getPayer(tx, payerID)
	if err != nil {
		return nil, nil, err
	}

	payee, err := svc.getPayee(tx, payeeID)
	if err != nil {
		return nil, nil, err
	}

	if payer.PayerID == payee.PayeeID {
		return nil, nil, errors.New("payer cannot pay themselves")
	}

	return payer, payee, nil
}

func (svc *TransactionService) getPayer(tx *gorm.DB, payerID string) (*model.Payer, error) {
	var payer model.Payer
	if err := tx.First(&payer, "PayerID = ?", payerID).Error; err != nil {
		errorLogger.Printf("Failed to retrieve payer with PayerID %s: %v\n", payerID, err)
		return nil, fmt.Errorf("failed to retrieve payer with PayerID %s: %v", payerID, err)
	}
	return &payer, nil
}

func (svc *TransactionService) getPayee(tx *gorm.DB, payeeID string) (*model.Payee, error) {
	var payee model.Payee
	if err := tx.First(&payee, "PayeeID = ?", payeeID).Error; err != nil {
		errorLogger.Printf("Failed to retrieve payee with PayeeID %s: %v\n", payeeID, err)
		return nil, fmt.Errorf("failed to retrieve payee with PayeeID %s: %v", payeeID, err)
	}
	return &payee, nil
}
func (svc *TransactionService) validatePaymentMethod(ctx context.Context, payerID, paymentMethodID string, paymentDetail model.PaymentDetails) error {
	// Fetch the payment method by payer ID and payment method ID
	paymentMethod, err := svc.GetPaymentMethodByPayerID(ctx, payerID, paymentMethodID)
	if err != nil {
		return err
	}

	// Check if the payment method is active
	if paymentMethod.Status != "active" {
		return errors.New("payment method is not active")
	}

	// Validate payment details (card number, expiry date)
	if err := svc.ValidatePaymentDetails(paymentMethod, paymentDetail); err != nil {
		return fmt.Errorf("invalid payment details: %v", err)
	}

	return nil
}
func (svc *TransactionService) GetPaymentMethodByPayerID(ctx context.Context, payerID, paymentMethodID string) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod

	if err := svc.DB.Where("payer_id = ? AND payment_method_id = ?", payerID, paymentMethodID).First(&paymentMethod).Error; err != nil {
		errorLogger.Printf("Failed to retrieve payment method for payerID %s with paymentMethodID %s: %v\n", payerID, paymentMethodID, err)
		return nil, fmt.Errorf("no valid payment method found for payer %s with paymentMethodID %s: %v", payerID, paymentMethodID, err)
	}

	return &paymentMethod, nil
}

func (svc *TransactionService) ValidatePaymentDetails(paymentMethod *model.PaymentMethod, paymentDetail model.PaymentDetails) error {
	if paymentMethod.CardNumber != paymentDetail.CardNumber || paymentMethod.ExpiryDate != paymentDetail.ExpiryDate {
		errorLogger.Println("Payment details validation failed: card number or expiry date mismatch")
		return errors.New("invalid payment details")
	}
	return nil
}

// func (svc *TransactionService) isDuplicateTransaction(tx *gorm.DB, payerID, payeeID string, amount float64, transactionType string) bool {
// 	var count int64
// 	if err := tx.Model(&model.Transaction{}).Where("payer_id = ? AND payee_id = ? AND amount = ? AND transaction_type = ?", payerID, payeeID, amount, transactionType).Count(&count).Error; err != nil {
// 		errorLogger.Printf("Error checking duplicate transaction for payer %s and payee %s: %v\n", payerID, payeeID, err)
// 		return false // Assume no duplicate if there's an error
// 	}
// 	return count > 0
// }

func (svc *TransactionService) hasSufficientBalance(payer *model.Payer, amount float64) bool {
	if payer.Balance < amount {
		errorLogger.Printf("Insufficient balance for payer %s. Available: %.2f, Required: %.2f\n", payer.PayerID, payer.Balance, amount)
		return false
	}
	return true
}

func (svc *TransactionService) createTransaction(tx *gorm.DB, payer *model.Payer, payee *model.Payee, amount float64, transactionType, status string, reservedAmount float64, paymentMethodID string, paymentDetail model.PaymentDetails) (*model.Transaction, error) {
	transactionID := utils.GenerateUniqueID()
	transaction := &model.Transaction{
		TransactionID:   transactionID,
		PayerID:         payer.PayerID,
		PayeeID:         payee.PayeeID,
		Amount:          amount,
		TransactionType: transactionType,
		Status:          status,
		ReservedAmount:  reservedAmount,
		PaymentMethodID: paymentMethodID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := tx.Create(transaction).Error; err != nil {
		errorLogger.Printf("Failed to create transaction %s: %v\n", transactionID, err)
		return nil, err
	}

	return transaction, nil
}

func (svc *TransactionService) processPayment(tx *gorm.DB, transaction *model.Transaction) error {
	if err := svc.debitAccount(tx, transaction.PayerID, transaction.Amount); err != nil {
		if err := svc.logFailedTransaction(tx, transaction.TransactionID, "Transaction Failed", "Failed to debit payer"); err != nil {
			errorLogger.Printf("Failed to log failed transaction audit in debitAccount: %v\n", err)
		}
		errorLogger.Printf("Failed to debit payer %s for transaction %s: %v\n", transaction.PayerID, transaction.TransactionID, err)
		return err
	}

	if err := svc.creditAccount(tx, transaction.PayeeID, transaction.Amount); err != nil {
		if err := svc.logFailedTransaction(tx, transaction.TransactionID, "Transaction Failed", "Failed to credit payee"); err != nil {
			errorLogger.Printf("Failed to log failed transaction audit in creditAccount: %v\n", err)
		}
		errorLogger.Printf("Failed to credit payee %s for transaction %s: %v\n", transaction.PayeeID, transaction.TransactionID, err)
		return err
	}

	transaction.Status = "completed"
	if err := tx.Save(transaction).Error; err != nil {
		if err := svc.logFailedTransaction(tx, transaction.TransactionID, "Transaction Failed", "Failed to update transaction"); err != nil {
			errorLogger.Printf("Failed to log failed transaction audit in completed state: %v\n", err)
		}
		errorLogger.Printf("Failed to update transaction %s to completed: %v\n", transaction.TransactionID, err)
		return err
	}

	return nil
}

func (svc *TransactionService) debitAccount(tx *gorm.DB, payerID string, amount float64) error {
	if err := tx.Model(&model.Payer{}).Where("PayerID = ?", payerID).Update("Balance", gorm.Expr("Balance - ?", amount)).Error; err != nil {
		errorLogger.Printf("Failed to debit account for payer %s: %v\n", payerID, err)
		return fmt.Errorf("failed to debit payer's account: %v", err)
	}
	return nil
}

func (svc *TransactionService) creditAccount(tx *gorm.DB, payeeID string, amount float64) error {
	if err := tx.Model(&model.Payee{}).Where("PayeeId = ?", payeeID).Update("Balance", gorm.Expr("Balance + ?", amount)).Error; err != nil {
		errorLogger.Printf("Failed to credit account for payee %s: %v\n", payeeID, err)
		return fmt.Errorf("failed to credit payee's account: %v", err)
	}
	return nil
}

func (svc *TransactionService) validateTransactionPayload(tx *gorm.DB, transactionID, payerID, payeeID string, amount float64, transactionType, paymentMethodID string) error {
	transactionType = strings.ToLower(transactionType)

	validTransactionTypes := map[string]bool{
		"debit":  true,
		"credit": true,
		"refund": true,
	}

	// Step 1: Validate transaction type
	if !validTransactionTypes[transactionType] {
		errorLogger.Printf("Invalid transaction type: %s\n", transactionType)
		return errors.New("invalid transaction type")
	}

	// Step 2: Validate amount (should be greater than 0)
	if amount <= 0 {
		errorLogger.Printf("Invalid transaction amount: %.2f for type: %s\n", amount, transactionType)
		return errors.New("transaction amount must be greater than zero")
	}

	// Step 3: If refund, ensure the amount matches the original transaction
	if transactionType == "refund" {
		var originalTransaction model.Transaction
		err := tx.First(&originalTransaction, "transaction_id = ?", transactionID).Error
		if err != nil {
			if err := svc.logFailedTransaction(tx, transactionID, "Transaction Failed", "original transaction not found for refund"); err != nil {
				errorLogger.Printf("Failed to log failed transaction audit in debitAccount: %v\n", err)
			}
			errorLogger.Printf("Original transaction %s not found for refund: %v\n", transactionID, err)
			return errors.New("original transaction not found for refund")
		}
		if originalTransaction.Amount != amount {
			if err := svc.logFailedTransaction(tx, transactionID, "Transaction Failed", "Refund amount mismatch"); err != nil {
				errorLogger.Printf("Failed to log failed transaction audit in debitAccount: %v\n", err)
			}
			errorLogger.Printf("Refund amount mismatch: original=%.2f, requested=%.2f\n", originalTransaction.Amount, amount)
			return errors.New("refund amount must match the original transaction amount")
		}
	}

	return nil
}

func (svc *TransactionService) ListTransactions(ctx context.Context, userID string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if err := svc.DB.Where("payer_id = ? OR payee_id = ?", userID, userID).Find(&transactions).Error; err != nil {
		errorLogger.Printf("Failed to fetch transactions for user %s: %v\n", userID, err)
		return nil, fmt.Errorf("failed to fetch transactions: %v", err)
	}
	return transactions, nil
}

func (svc *TransactionService) createAuditEntry(tx *gorm.DB, transaction *model.Transaction) error {
	fmt.Println("Transaction details", transaction)
	err := svc.logAudit(tx, transaction.TransactionID, "Transaction Created", "Transaction successfully created and processed")
	if err != nil {
		errorLogger.Printf("Failed to create audit entry for transaction %s: %v\n", transaction.TransactionID, err)
	}
	return err
}

func (svc *TransactionService) logFailedTransaction(tx *gorm.DB, transactionID, action, details string) error {
	auditLog := model.AuditLog{
		AuditLogID:    utils.GenerateUniqueID(),
		TransactionID: transactionID,
		Action:        action,
		Details:       details,
		CreatedAt:     time.Now(),
	}

	// Log failed transaction audit
	if err := tx.Create(&auditLog).Error; err != nil {
		errorLogger.Printf("Failed to log failed audit for transaction %s: %v\n", transactionID, err)
		return err
	}

	return nil
}

func (svc *TransactionService) logAudit(tx *gorm.DB, transactionID, action, details string) error {
	auditLog := model.AuditLog{
		AuditLogID:    utils.GenerateUniqueID(),
		TransactionID: transactionID,
		Action:        action,
		Details:       details,
		CreatedAt:     time.Now(),
	}

	if err := tx.Create(&auditLog).Error; err != nil {
		errorLogger.Printf("Failed to log audit for transaction %s: %v\n", transactionID, err)
		return err
	}

	return nil
}
