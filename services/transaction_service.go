package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"poc/model"
	"poc/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

var errorLogger *log.Logger

func init() {
	errorLogFile, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
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

func (svc *TransactionService) getPayerAndPayee(payerID, payeeID string) (*model.Payer, *model.Payee, error) {
	// Retrieve payer
	var payer model.Payer
	if err := svc.DB.First(&payer, "PayerID = ?", payerID).Error; err != nil {
		// Log the error to the file
		errorLogger.Printf("Failed to retrieve payer with PayerID %s: %v\n", payerID, err)
		// Return nil for payer and a custom error
		return nil, nil, fmt.Errorf("failed to retrieve payer with PayerID %s: %v", payerID, err)
	}

	// Retrieve payee
	var payee model.Payee
	if err := svc.DB.First(&payee, "PayeeID = ?", payeeID).Error; err != nil {
		// Log the error to the file
		errorLogger.Printf("Failed to retrieve payee with PayeeID %s: %v\n", payeeID, err)
		// Return nil for payee and a custom error
		return nil, nil, fmt.Errorf("failed to retrieve payee with PayeeID %s: %v", payeeID, err)
	}

	// Return payer and payee if no errors occurred
	return &payer, &payee, nil
}

func (svc *TransactionService) InitializeTransaction(ctx context.Context, payerID, payeeID string, amount float64, transactionType, status string, reservedAmount float64, transaction_id, paymentMethodID string, paymentDetail model.PaymentDetails) (*model.Transaction, error) {

	var payer *model.Payer
	var payee *model.Payee
	payer, payee, err := svc.getPayerAndPayee(payerID, payeeID)
	if err != nil {
		// Handle the error (already logged to file)
		fmt.Println("Error:", err)
		return nil, err
	}
	// Step 3: Validate that payer and payee are not the same
	if payer.PayerID == payee.PayeeID {
		return nil, errors.New("payer cannot pay themselves")
	}

	// Step 4: Fetch and validate the payment method

	// Main code for calling the GetPaymentMethodByPayerAndPaymentMethodID
	getPaymentMethod, err := svc.GetPaymentMethodByPayerAndPaymentMethodID(ctx, payer.PayerID, paymentMethodID, transactionType == "Refund") // true indicates it's a refund
	if err != nil {
		// Log the error and handle the refund case
		errorLogger.Printf("Error retrieving payment method for refund for PayerID %s and PaymentMethodID %s: %v\n", payer.PayerID, paymentMethodID, err)
		// Additional handling logic here, like returning the error to the caller or a custom response
		return nil, fmt.Errorf("payment method GetPaymentMethodByPayerAndPaymentMethodID failed: %v", err)
	}

	// Pass 'true' for isRefund to skip the active status check
	err = svc.ValidatePaymentDetails(getPaymentMethod, paymentDetail, false) // transactionType == "Refund")
	if err != nil {
		// Log error before returning
		errorLogger.Printf("Payment method validation failed for PayerID %s, PaymentMethodID %s: %v\n", payerID, paymentMethodID, err)
		return nil, fmt.Errorf("payment method validation failed: %v", err)
	}
	transactionID := utils.GenerateUniqueID()
	// Step 5: Validate the transaction payload
	// if transactionType == "Refund" {
	// 	transactionID = transaction_id
	// }
	if err := validateTransactionPayload(transactionID, payerID, payeeID, amount, transactionType, paymentMethodID); err != nil {
		return nil, fmt.Errorf("invalid transaction payload: %v", err)
	}
	if transactionType != "Refund" {
		// Step 6: Check for duplicate transaction
		if err := svc.CheckDuplicateTransaction(transactionID); err != nil {
			return nil, fmt.Errorf("duplicate transaction: %v", err)
		}
	}

	// Step 7: Create the transaction record
	transaction := &model.Transaction{
		TransactionID:   transactionID,
		PayerID:         payerID,
		PayeeID:         payeeID,
		Amount:          amount,
		TransactionType: transactionType,
		Status:          "Pending",
		ReservedAmount:  reservedAmount,
		PaymentMethodID: paymentMethodID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	defer func() {
		auditLog := model.AuditLog{
			AuditLogID:    utils.GenerateUniqueID(),
			TransactionID: transactionID,
			Action:        "Transaction Process Failed",
			Details:       "Failed",
			CreatedAt:     time.Now(),
		}
		if err := svc.DB.Create(&auditLog).Error; err != nil {
			fmt.Printf("Failed to log audit entry: %v\n", err)
		}
	}()
	if err := svc.DB.Create(transaction).Error; err != nil {
		errorLogger.Printf("Failed Create Transaction : %v\n", err)
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	if transaction.TransactionType == "Debit" {
		// Step 8: Check the payer's balance
		if err := svc.CheckBalance(ctx, transaction); err != nil {
			errorLogger.Printf("Balance check failed : %v\n", err)
			return nil, fmt.Errorf("balance check failed: %v", err)
		}
		// Step 9: Reserve the funds
		if err := svc.ReserveFunds(ctx, transaction); err != nil {
			errorLogger.Printf("Failed reserve funds : %v\n", err)
			return nil, fmt.Errorf("failed to reserve funds: %v", err)
		}
	}

	// Step 10: Process the payment
	if err := svc.ProcessPayment(ctx, transaction); err != nil {
		errorLogger.Printf("Process Payment failed : %v\n", err)
		return nil, fmt.Errorf("Payment processing failed: %v", err)
	}

	// Step 11: Complete the transaction
	if err := svc.CompleteTransaction(transaction.TransactionID, svc.DB); err != nil {
		errorLogger.Printf("Failed to complete transaction : %v\n", err)
		return nil, fmt.Errorf("failed to complete transaction: %v", err)
	}

	return transaction, nil
}
func (svc *TransactionService) ValidatePaymentDetails(paymentMethod *model.PaymentMethod, paymentDetail model.PaymentDetails, isRefund bool) error {
	// If it's not a refund, we check if the payment method is active
	if !isRefund && paymentMethod.Status != "active" {
		// Log error before returning
		errorLogger.Printf("Payment method with ID %s is not active. Payment method status: %s\n", paymentMethod.PaymentMethodID, paymentMethod.Status)
		return errors.New("payment method is not active")
	}

	switch paymentMethod.MethodType {
	case "card":
		if paymentMethod.CardNumber != paymentDetail.CardNumber {
			// Log error before returning
			errorLogger.Printf("Card number mismatch for payment method ID %s. Provided: %s, Stored: %s\n", paymentMethod.PaymentMethodID, paymentDetail.CardNumber, paymentMethod.CardNumber)
			return errors.New("payment method is not correct - card number")
		}
		if paymentMethod.ExpiryDate != paymentDetail.ExpiryDate {
			// Log error before returning
			errorLogger.Printf("Expiry date mismatch for payment method ID %s. Provided: %s, Stored: %s\n", paymentMethod.PaymentMethodID, paymentDetail.ExpiryDate, paymentMethod.ExpiryDate)
			return errors.New("payment method is not correct - expiry date")
		}
	case "bank_transfer":
		if paymentMethod.AccountNumber != paymentDetail.AccountNumber {
			// Log error before returning
			errorLogger.Printf("Account number mismatch for payment method ID %s. Provided: %s, Stored: %s\n", paymentMethod.PaymentMethodID, paymentDetail.AccountNumber, paymentMethod.AccountNumber)
			return errors.New("payment method is bank_transfer, please provide correct - account number")
		}
	case "upi":
		if paymentMethod.Details != paymentDetail.UPIID {
			// Log error before returning
			errorLogger.Printf("UPI ID mismatch for payment method ID %s. Provided: %s, Stored: %s\n", paymentMethod.PaymentMethodID, paymentDetail.UPIID, paymentMethod.Details)
			return errors.New("payment method is not correct - upi id")
		}
	case "wallet":
		if paymentMethod.Details != paymentDetail.Wallet {
			// Log error before returning
			errorLogger.Printf("Wallet mismatch for payment method ID %s. Provided: %s, Stored: %s\n", paymentMethod.PaymentMethodID, paymentDetail.Wallet, paymentMethod.Details)
			return errors.New("payment method is not correct - wallet")
		}
	case "cheque":
		if paymentMethod.Details != paymentDetail.Cheque {
			// Log error before returning
			errorLogger.Printf("Cheque mismatch for payment method ID %s. Provided: %s, Stored: %s\n", paymentMethod.PaymentMethodID, paymentDetail.Cheque, paymentMethod.Details)
			return errors.New("payment method is not correct - cheque")
		}
	default:
		// Log error for invalid payment method
		errorLogger.Printf("Invalid payment method type for payment method ID %s: %s\n", paymentMethod.PaymentMethodID, paymentMethod.MethodType)
		return errors.New("invalid payment method type")
	}

	return nil
}

func (svc *TransactionService) GetPaymentMethodByPayerAndPaymentMethodID(ctx context.Context, payerID, paymentMethodID string, isRefund bool) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	var err error

	// If it's a refund, skip the "active" status check
	if isRefund {
		fmt.Println("Camer here at refund")
		// Retrieve payment method for a refund (ignoring active status check)
		if err = svc.DB.First(&paymentMethod, "payer_id = ? AND payment_method_id = ?", payerID, paymentMethodID).Error; err != nil {
			// Log the error and return a custom error message
			errorLogger.Printf("Failed to retrieve payment method for PayerID %s, PaymentMethodID %s (Refund Case): %v\n", payerID, paymentMethodID, err)
			return nil, fmt.Errorf("failed to retrieve payment method details for PayerID %s, PaymentMethodID %s: %v", payerID, paymentMethodID, err)
		}
	} else {
		fmt.Println("Camer here at other")
		// Retrieve payment method with active status check
		err = svc.DB.Where("payer_id = ? AND payment_method_id = ? AND status = ?", payerID, paymentMethodID, "active").First(&paymentMethod).Error
		if err != nil {
			// Log the error and return a custom error message
			errorLogger.Printf("Failed to retrieve active payment method for PayerID %s, PaymentMethodID %s: %v\n", payerID, paymentMethodID, err)
			return nil, fmt.Errorf("no valid payment method found for payer %s with payment method %s: %v", payerID, paymentMethodID, err)
		}
	}

	// Return the payment method details if successful
	fmt.Println("payment method: ", &paymentMethod)
	return &paymentMethod, nil
}

func validateTransactionPayload(transactionID, payerID, payeeID string, amount float64, transactionType, paymentMethodID string) error {
	transactionType = strings.ToLower(transactionType)

	validTransactionTypes := map[string]bool{
		"debit":  true,
		"credit": true,
		"refund": true,
	}

	if !validTransactionTypes[transactionType] {
		errorLogger.Printf("Invalid transaction type: %s\n", transactionType)
		return errors.New("invalid transaction type")
	}

	// Validate required fields for "Debit" transactions
	if (transactionType == "debit") && (transactionID == "" || payerID == "" || payeeID == "" || amount <= 0 || paymentMethodID == "") {
		errorLogger.Printf("Invalid input type: %s\n %s\n %s\n %s\n %s\n %s\n%s\n", transactionType, transactionID, payerID, payeeID, amount, paymentMethodID)
		return errors.New("missing required fields or invalid data for Debit transaction")
	}

	return nil
}

func (svc *TransactionService) CheckDuplicateTransaction(transactionID string) error {
	var existingTransaction model.Transaction
	if err := svc.DB.First(&existingTransaction, "transaction_id = ?", transactionID).Error; err == nil {
		errorLogger.Printf("Duplicate transaction ID check failed : %v\n", err)
		return errors.New("duplicate transaction ID")
	}
	return nil
}

func (svc *TransactionService) CheckBalance(ctx context.Context, transaction *model.Transaction) error {
	var payer model.Payer
	if err := svc.DB.First(&payer, "PayerID = ?", transaction.PayerID).Error; err != nil {
		errorLogger.Printf("Payer not found: %s\n", transaction.PayerID)
		return errors.New("payer not found")
	}
	if payer.Balance < transaction.Amount {
		errorLogger.Printf("Payer Balance is less then Transaction amount : %s\n %s\n", payer.Balance, transaction.Amount)
		_ = svc.UpdateTransactionStatus(ctx, transaction.TransactionID, "Failed")
		return errors.New("insufficient funds")
	}
	return nil
}
func (svc *TransactionService) ReserveFunds(ctx context.Context, transaction *model.Transaction) error {
	return svc.DB.Transaction(func(tx *gorm.DB) error {
		var payer model.Payer
		if err := tx.First(&payer, "PayerID = ?", transaction.PayerID).Error; err != nil {
			return errors.New("payer not found")
		}
		if payer.Balance < transaction.Amount {
			_ = svc.UpdateTransactionStatus(ctx, transaction.TransactionID, "Failed")
			return errors.New("insufficient funds")
		}
		payer.Balance -= transaction.Amount
		if err := tx.Save(&payer).Error; err != nil {
			return err
		}
		fmt.Println("transaction - Status", transaction.Status)
		transaction.Status = "Reserved"
		transaction.ReservedAmount = transaction.Amount
		return tx.Save(transaction).Error
	})
}
func (svc *TransactionService) RollbackReservation(ctx context.Context, transaction *model.Transaction) error {
	return svc.DB.Transaction(func(tx *gorm.DB) error {
		var payer model.Payer
		if err := tx.First(&payer, "PayerID = ?", transaction.PayerID).Error; err != nil {
			return errors.New("payer not found")
		}

		// Add back the reserved amount
		payer.Balance += transaction.ReservedAmount
		if err := tx.Save(&payer).Error; err != nil {
			return fmt.Errorf("failed to rollback reservation: %v", err)
		}

		transaction.Status = "Failed"
		transaction.ReservedAmount = 0
		return tx.Save(transaction).Error
	})
}

func (svc *TransactionService) ProcessPayment(ctx context.Context, transaction *model.Transaction) error {
	return svc.DB.Transaction(func(tx *gorm.DB) error {
		// Ensure the reserved funds are rolled back on failure
		defer func() {
			if r := recover(); r != nil {
				transaction.Status = "Failed"
				svc.RollbackReservation(ctx, transaction)
			}
		}()

		switch transaction.TransactionType {
		case "Debit":
			var payer model.Payer
			var payee model.Payee

			// // Fetch payer details
			if err := tx.First(&payer, "PayerID = ?", transaction.PayerID).Error; err != nil {
				return errors.New("payer not found")
			}

			// Fetch payee details
			if err := tx.First(&payee, "PayeeID = ?", transaction.PayeeID).Error; err != nil {
				return errors.New("payee not found")
			}

			payee.Balance += transaction.Amount

			// Save updated records
			if err := tx.Save(&payer).Error; err != nil {
				return err
			}
			if err := tx.Save(&payee).Error; err != nil {
				return err
			}

		case "Credit":
			var payee model.Payee
			if err := tx.First(&payee, "PayeeID = ?", transaction.PayeeID).Error; err != nil {
				return errors.New("payee not found")
			}

			// Update balance
			payee.Balance += transaction.Amount
			if err := tx.Save(&payee).Error; err != nil {
				return err
			}

		case "Refund":
			// Validate original transaction
			var originalTransaction model.Transaction
			if err := tx.First(&originalTransaction, "transaction_id = ?", transaction.TransactionID).Error; err != nil {
				return errors.New("original transaction not found")
			}
			// Fetch payer and payee from the original transaction
			var payer model.Payer
			var payee model.Payee

			if err := tx.First(&payer, "PayerID = ?", originalTransaction.PayerID).Error; err != nil {
				return errors.New("payer not found")
			}
			if err := tx.First(&payee, "PayeeID = ?", originalTransaction.PayeeID).Error; err != nil {
				return errors.New("payee not found")
			}

			// Reverse balances
			if payee.Balance <= 0 && payee.Balance < originalTransaction.Amount {
				return errors.New("insufficient funds in payee account for refund")
			}

			payer.Balance += originalTransaction.Amount
			payee.Balance -= originalTransaction.Amount

			// Save updated records
			if err := tx.Save(&payer).Error; err != nil {
				return err
			}
			if err := tx.Save(&payee).Error; err != nil {
				return err
			}

			// Mark original transaction as refunded
			originalTransaction.Status = "Refunded"
			if err := tx.Save(&originalTransaction).Error; err != nil {
				return err
			}

		default:
			return errors.New("unsupported transaction type")
		}

		// Mark transaction as completed
		transaction.ReservedAmount = 0
		transaction.Status = "Completed"
		return tx.Save(transaction).Error
	})
}
func RefundTransaction(transactionID string, db *gorm.DB) error {
	// Start a database transaction
	return db.Transaction(func(tx *gorm.DB) error {
		// Fetch the transaction
		var transaction model.Transaction
		if err := tx.First(&transaction, "transaction_id = ?", transactionID).Error; err != nil {
			return fmt.Errorf("transaction not found: %w", err)
		}

		// Check if the transaction is already refunded
		if transaction.Status == "Refunded" {
			return fmt.Errorf("transaction already refunded")
		}

		// Fetch the payee
		var payee model.Payee
		if err := tx.First(&payee, "payee_id = ?", transaction.PayeeID).Error; err != nil {
			return fmt.Errorf("payee not found: %w", err)
		}

		// Fetch the payer
		var payer model.Payer
		if err := tx.First(&payer, "payer_id = ?", transaction.PayerID).Error; err != nil {
			return fmt.Errorf("payer not found: %w", err)
		}

		// Check if the payee has sufficient balance for the refund
		if payee.Balance < transaction.Amount {
			return fmt.Errorf("insufficient balance in payee account for refund")
		}

		// Perform the refund by adjusting balances
		payee.Balance -= transaction.Amount
		payer.Balance += transaction.Amount

		// Save updated balances
		if err := tx.Save(&payee).Error; err != nil {
			return fmt.Errorf("failed to update payee balance: %w", err)
		}

		if err := tx.Save(&payer).Error; err != nil {
			return fmt.Errorf("failed to update payer balance: %w", err)
		}

		// Update the transaction status to refunded
		if err := tx.Model(&transaction).Update("status", "Refunded").Error; err != nil {
			return fmt.Errorf("error updating transaction status to refunded: %w", err)
		}

		// Add entry to the audit log
		auditLog := model.AuditLog{
			AuditLogID:    utils.GenerateUniqueID(),
			TransactionID: transactionID,
			Action:        "Transaction success",
			Details:       "Transaction refunded successfully",
			CreatedAt:     time.Now(),
		}
		if err := tx.Create(&auditLog).Error; err != nil {
			return fmt.Errorf("error inserting audit log: %w", err)
		}

		return nil
	})
}

func (svc *TransactionService) CompleteTransaction(transactionID string, db *gorm.DB) error {
	// Update transaction status
	if err := db.Model(&model.Transaction{}).Where("transaction_id = ?", transactionID).
		Update("status", "Completed").Error; err != nil {
		return fmt.Errorf("error completing transaction: %w", err)
	}

	// Add entry to audit log
	// auditLog := model.AuditLog{
	// 	AuditLogID:    utils.GenerateUniqueID(),
	// 	TransactionID: transactionID,
	// 	Action:        "Transaction Success",
	// 	Details:       "Transaction completed successfully",
	// 	CreatedAt:     time.Now(),
	// }
	// svc.logAudit(transactionID, "Transaction Success", "Transaction completed successfully")
	// if err := db.Create(&auditLog).Error; err != nil {
	// 	return fmt.Errorf("error inserting audit log: %w", err)
	// }

	return nil
}

func (svc *TransactionService) GetTransactionByID(ctx context.Context, transactionID string) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := svc.DB.Preload("Payer").Preload("Payee").Preload("PaymentMethod").First(&transaction, "transaction_id = ?", transactionID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %v", err)
	}
	return &transaction, nil
}

// ListTransactions retrieves all transactions for a specific user as payer or payee
func (svc *TransactionService) ListTransactions(ctx context.Context, userID string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if err := svc.DB.Where("payer_id = ? OR payee_id = ?", userID, userID).Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %v", err)
	}
	return transactions, nil
}

// UpdateTransactionStatus updates the status of a transaction
func (svc *TransactionService) UpdateTransactionStatus(ctx context.Context, transactionID, status string) error {
	if err := svc.DB.Model(&model.Transaction{}).Where("transaction_id = ?", transactionID).Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update transaction status: %v", err)
	}
	return nil
}

// DeleteTransaction deletes a transaction by its ID
func (svc *TransactionService) DeleteTransaction(ctx context.Context, transactionID string) error {
	if err := svc.DB.Delete(&model.Transaction{}, "transaction_id = ?", transactionID).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %v", err)
	}
	return nil
}

func (svc *TransactionService) logAudit(transactionID, action, details string) {
	auditLog := model.AuditLog{
		AuditLogID:    utils.GenerateUniqueID(),
		TransactionID: transactionID,
		Action:        action,
		Details:       details,
		CreatedAt:     time.Now(),
	}

	if err := svc.DB.Create(&auditLog).Error; err != nil {
		log.Printf("Failed to log audit: %v", err)
	}
}
