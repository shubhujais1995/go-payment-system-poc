package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"poc/model"
	"poc/utils"
	"time"

	"gorm.io/gorm"
)

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

func (svc *TransactionService) InitializeTransaction(ctx context.Context, payerID, payeeID string, amount float64, transactionType, status string, reservedAmount float64, paymentMethodID string, paymentDetail model.PaymentDetails) (*model.Transaction, error) {

	// err := svc.DB.Transaction(func(tx *gorm.DB) error {

	// 	var payer model.Payer
	// 	if err := svc.DB.First(&payer, "PayerID = ?", payerID).Error; err != nil {
	// 		return err
	// 	}

	// 	// Step 3: Check if the payee exists
	// 	var payee model.Payee
	// 	if err := svc.DB.First(&payee, "PayeeID = ?", payeeID).Error; err != nil {
	// 		return err
	// 	}

	// 	// return nil will commit the whole transaction
	// 	return nil
	// })

	/***********************/

	// Step 1: Check if the payer exists
	var payer model.Payer
	if err := svc.DB.First(&payer, "PayerID = ?", payerID).Error; err != nil {
		return nil, fmt.Errorf("payer with PayerID %s does not exist", payerID)
	}

	// if transactionType == "Debit" {
	// Step 2: Fetch and validate the payment method
	// paymentMethod, err := svc.GetPaymentMethodByPayerID(ctx, payerID)
	paymentMethod, err := svc.GetPaymentMethodByPayerAndPaymentMethodID(ctx, payerID, paymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("no valid payment method found for payer GetPaymentMethodByPayerID: %v", err)
	}

	if paymentMethod.Status != "active" {
		return nil, errors.New("payment method is not active")
	}
	// }

	fmt.Println("Reached Here!", "paymentDetail.AccountNumber -", paymentDetail.AccountNumber, "paymentDetail.CardNumber=", paymentDetail.CardNumber,
		" paymentDetail.CVV=", paymentDetail.CVV, "paymentDetail.ExpiryDate=", paymentDetail.ExpiryDate)
	fmt.Println("paymentMethod came from server", paymentMethod.CardNumber, paymentMethod.ExpiryDate, paymentMethod.MethodType, "paymentMethod.PaymentMethodID -", paymentMethod.PaymentMethodID)
	// if transactionType == "Debit" {
	errPaymentMethod := svc.ValidatePaymentDetails(paymentMethod, paymentDetail)
	if errPaymentMethod != nil {
		return nil, fmt.Errorf("no valid payment method found for payer ValidatePaymentDetails: %v", errPaymentMethod)
	}
	// }

	// Step 3: Check if the payee exists
	var payee model.Payee
	if err := svc.DB.First(&payee, "PayeeID = ?", payeeID).Error; err != nil {
		return nil, fmt.Errorf("payee with PayeeID %s does not exist", payeeID)
	}

	// Step 4: Validate that payer and payee are not the same
	if payerID == payeeID {
		return nil, errors.New("payer cannot pay themselves")
	}

	// Step 5: Validate the transaction payload
	transactionID := utils.GenerateUniqueID()
	if err := validateTransactionPayload(transactionID, payerID, payeeID, amount, transactionType, paymentMethodID); err != nil {
		return nil, fmt.Errorf("invalid transaction payload: %v", err)
	}

	// Step 6: Check for duplicate transaction
	if err := svc.CheckDuplicateTransaction(transactionID); err != nil {
		return nil, fmt.Errorf("duplicate transaction: %v", err)
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
		// svc.logAudit(transaction.TransactionID, "Create T", "Trasaction failed")
		return nil, fmt.Errorf("failed to create transaction: %v", err)
	}

	if transactionType == "Debit" {
		// Step 8: Check the payer's balance
		if err := svc.CheckBalance(ctx, transaction); err != nil {
			return nil, fmt.Errorf("balance check failed: %v", err)
		}
	}

	// // Step 9: Reserve the funds
	// if err := svc.ReserveFunds(ctx, transaction); err != nil {
	// 	return nil, fmt.Errorf("failed to reserve funds: %v", err)
	// }

	if transaction.TransactionType == "Debit" {
		// 	// Skip reserving funds for refunds
		// 	// You can directly proceed to reverse balances
		// Step 9: Reserve the funds
		if err := svc.ReserveFunds(ctx, transaction); err != nil {
			return nil, fmt.Errorf("failed to reserve funds: %v", err)
		}
	}

	// Step 10: Process the payment
	if err := svc.ProcessPayment(ctx, transaction); err != nil {

		return nil, fmt.Errorf("payment processing failed: %v", err)
	}

	// Step 11: Complete the transaction
	if err := svc.CompleteTransaction(transaction.TransactionID, svc.DB); err != nil {

		return nil, fmt.Errorf("failed to complete transaction: %v", err)
	}

	return transaction, nil
}

func (svc *TransactionService) ValidatePaymentDetails(paymentMethod *model.PaymentMethod, paymentDetail model.PaymentDetails) error {
	switch paymentMethod.MethodType {
	case "card":

		if paymentMethod.CardNumber != paymentDetail.CardNumber {
			return errors.New("payment method is not correct - card number")
		}
		// if paymentMethod.cvv != paymentDetail.PaymentDetails.CVV {
		// 	return nil, errors.New("payment method is not correct cvv")
		// }
		if paymentMethod.ExpiryDate != paymentDetail.ExpiryDate {
			return errors.New("payment method is not correct - expiry date")
		}
	case "bank_transfer":
		// Validate AccountNumber for bank transfer
		if paymentMethod.AccountNumber != paymentDetail.AccountNumber {
			fmt.Println(paymentMethod.AccountNumber, "paymentMethod.AccountNumber", paymentDetail.AccountNumber, "paymentDetail.AccountNumber")
			return errors.New("payment method is bank_transfer, please provide correct - account number")
		}

	case "upi":
		// Validate UPI if required
		if paymentMethod.Details != paymentDetail.UPIID {
			return errors.New("payment method is not correct - upi id")
		}

	case "wallet":
		// Validate Wallet method if required
		if paymentMethod.Details != paymentDetail.Wallet {
			return errors.New("payment method is not correct - wallet")
		}

	case "cheque":
		// Validate for cheque method if required
		if paymentMethod.Details != paymentDetail.Cheque {
			return errors.New("payment method is not correct - cheque")
		}

	default:
		return errors.New("invalid payment method type")

	}

	return nil
}
func (svc *TransactionService) UpdatePayerBalance(ctx context.Context, payerID string, amount float64) error {
	// Validate the amount to prevent invalid updates
	if amount == 0 {
		return errors.New("amount must not be zero")
	}

	// Begin a database transaction
	return svc.DB.Transaction(func(tx *gorm.DB) error {
		// Fetch the payer record
		var payer model.Payer
		if err := tx.First(&payer, "payer_id = ?", payerID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("payer with ID %s not found", payerID)
			}
			return fmt.Errorf("error fetching payer: %v", err)
		}

		// Update the payer's balance
		newBalance := payer.Balance + amount
		if newBalance < 0 {
			return errors.New("insufficient funds for balance update")
		}

		payer.Balance = newBalance

		// Save the updated balance
		if err := tx.Save(&payer).Error; err != nil {
			return fmt.Errorf("failed to update payer's balance: %v", err)
		}

		// Optionally log the balance update
		// logEntry := model.AuditLog{
		// 	LogID:      utils.GenerateUniqueID(),
		// 	PayerID:    payerID,
		// 	Amount:     amount,
		// 	NewBalance: newBalance,
		// 	CreatedAt:  time.Now(),
		// }
		// if err := tx.Create(&logEntry).Error; err != nil {
		// 	return fmt.Errorf("failed to log balance update: %v", err)
		// }

		return nil
	})
}

func (svc *TransactionService) DepositToPayer(ctx context.Context, payerID string, amount float64) error {
	// Validate the deposit amount
	if amount <= 0 {
		return errors.New("deposit amount must be greater than zero")
	}

	// Fetch the payer record
	var payer model.Payer
	if err := svc.DB.First(&payer, "payer_id = ?", payerID).Error; err != nil {
		return fmt.Errorf("payer with ID %s not found: %v", payerID, err)
	}

	// Create a deposit transaction
	depositTransaction := &model.Transaction{
		TransactionID:   utils.GenerateUniqueID(),
		PayerID:         payerID,
		PayeeID:         "", // No payee for deposit
		Amount:          amount,
		TransactionType: "Deposit",
		Status:          "Completed",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save the deposit transaction
	if err := svc.DB.Create(depositTransaction).Error; err != nil {
		return fmt.Errorf("failed to create deposit transaction: %v", err)
	}

	// Update the payer's balance
	payer.Balance += amount
	if err := svc.DB.Save(&payer).Error; err != nil {
		return fmt.Errorf("failed to update payer's balance: %v", err)
	}

	return nil
}

//	func (svc *TransactionService) GetPaymentMethodByPayerID(ctx context.Context, payerID string) (*model.PaymentMethod, error) {
//		var paymentMethod model.PaymentMethod
//		err := svc.DB.Where("payer_id = ? AND status = ?", payerID, "active").First(&paymentMethod).Error
//		if err != nil {
//			return nil, fmt.Errorf("no valid payment method found for payer GetPaymentMethodByPayerID  inside method: %v", err)
//		}
//		fmt.Println("payment method - 1 : ", &paymentMethod)
//		return &paymentMethod, nil
//	}
func (svc *TransactionService) GetPaymentMethodByPayerAndPaymentMethodID(ctx context.Context, payerID string, paymentMethodID string) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	err := svc.DB.Where("payer_id = ? AND payment_method_id = ? AND status = ?", payerID, paymentMethodID, "active").First(&paymentMethod).Error
	if err != nil {
		return nil, fmt.Errorf("no valid payment method found for payer %s with payment method %s: %v", payerID, paymentMethodID, err)
	}
	fmt.Println("payment method - 1 : ", &paymentMethod)
	return &paymentMethod, nil
}

func validateTransactionPayload(transaction_id string, payerID string, payeeID string, amount float64, transactionType, paymentMethodID string) error {
	// if transaction_id == "" || payerID == "" || payeeID == "" ||
	// 	amount <= 0 || transactionType == "" || paymentMethodID == "" {
	// 	return errors.New("missing required fields or invalid data")
	// }

	if transactionType != "Debit" && transactionType != "Credit" && transactionType != "Refund" && transactionType != "debit" && transactionType != "credit" && transactionType != "refund" {
		return errors.New("invalid transaction type")
	}
	if transactionType == "Debit" && (transaction_id == "" || payerID == "" || payeeID == "" || amount <= 0 || transactionType == "") {
		if paymentMethodID == "" {
			return errors.New("missing required fields or invalid data")
		}
	}
	return nil
}

func (svc *TransactionService) CheckDuplicateTransaction(transactionID string) error {
	var existingTransaction model.Transaction
	if err := svc.DB.First(&existingTransaction, "transaction_id = ?", transactionID).Error; err == nil {
		return errors.New("duplicate transaction ID")
	}
	return nil
}
func (svc *TransactionService) VerifyPaymentMethod(ctx context.Context, transaction *model.Transaction) error {
	paymentMethod, err := svc.PaymentMethodService.ValidatePaymentMethod(transaction.PaymentMethodID)
	if err != nil || paymentMethod.Status != "active" {
		_ = svc.UpdateTransactionStatus(ctx, transaction.TransactionID, "Failed")
		return errors.New("invalid or inactive payment method")
	}
	return nil
}
func (svc *TransactionService) CheckBalance(ctx context.Context, transaction *model.Transaction) error {
	var payer model.Payer
	if err := svc.DB.First(&payer, "PayerID = ?", transaction.PayerID).Error; err != nil {
		// svc.logAudit(transaction.TransactionID, "Check balance", "Trasaction failed")
		return errors.New("payer not found")
	}
	if payer.Balance < transaction.Amount {
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

			// // Check for sufficient balance
			// if payer.Balance < transaction.Amount {
			// 	return errors.New("insufficient funds")
			// }

			// Update balances
			// payer.Balance -= transaction.Amount
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

			fmt.Println("originalTransaction Status ", originalTransaction.Status, " - for - ", originalTransaction.TransactionID, "- Transaction type - ", originalTransaction.TransactionType)
			// Ensure the original transaction was completed
			// if originalTransaction.Status != "Completed" {
			// 	return errors.New("refund not allowed for incomplete transactions")
			// }

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
