package model

import "time"

type Transaction struct {
	TransactionID   string  `gorm:"primaryKey;size:36"`                    // Unique transaction identifier
	PayerID         string  `gorm:"not null;index"`                        // Foreign key to the Payer table
	Payer           Payer   `gorm:"foreignKey:PayerID;references:PayerID"` // Link to Payer details
	PayeeID         string  `gorm:"not null;index"`                        // Foreign key to the Payee table
	Payee           Payee   `gorm:"foreignKey:PayeeID;references:PayeeID"` // Link to Payee details
	Amount          float64 `gorm:"not null"`                              // Total transaction amount
	ReservedAmount  float64 `gorm:"default:0.0"`                           // Amount reserved, if any
	TransactionType string  `gorm:"size:20;not null"`                      // Type of transaction (Debit, Credit, Refund)
	Status          string  `gorm:"size:20;not null"`                      // Status of the transaction (Pending, Completed, Failed, Reserved)
	// Remove this if you do not want this dependency:
	PaymentMethodID string `gorm:"size:36;index"` // Foreign key to PaymentMethod table
	//PaymentMethod   PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:PaymentMethodID"` // Link to Payment Method details (remove if not needed)
	CreatedAt time.Time `gorm:"autoCreateTime"` // Timestamp for when the transaction was created
	UpdatedAt time.Time `gorm:"autoUpdateTime"` // Timestamp for when the transaction was last updated
}

// TableName explicitly sets the table name to "Transactions" (case-sensitive)
func (Transaction) TableName() string {
	return "Transactions" // Use the exact case-sensitive name you want
}

type PaymentDetails struct {
	AccountNumber string `json:"accountNumber"`
	CardNumber    string `json:"card_number" validate:"required"`
	CVV           string `json:"cvv" validate:"required"`
	ExpiryDate    string `json:"expiry_date" validate:"required"`
	UPIID         string `json:"upi_id" validate: "required"`
	Wallet        string `json:"wallet" validate: "required"`
	Cheque        string `json:"cheque" validate: "required"`
}

type ProcessPaymentInput struct {
	TransactionID   string         `json:"transaction_id" validate:"required"`
	PayerID         string         `json:"payer_id" validate:"required"`
	PayeeID         string         `json:"payee_id" validate:"required"`
	Status          string         `json:"status" validate:"required"`
	Amount          float64        `json:"amount" validate:"required,gt=0"`
	TransactionType string         `json:"transaction_type" validate:"required"`
	PaymentMethodID string         `json:"payment_method_id" validate:"required"`
	PaymentDetails  PaymentDetails `json:"payment_details"`
}
