package model

import "time"

// PaymentMethod represents a payment method associated with a payer.
type PaymentMethod struct {
	PaymentMethodID string    `gorm:"primaryKey;column:payment_method_id;size:36"` // Unique identifier for each payment method
	PayerID         string    `gorm:"not null;index"`                              // Foreign key referencing the payer
	MethodType      string    `gorm:"size:20;not null"`                            // Type of payment method (e.g., card, bank_transfer, wallet)
	CardNumber      string    `gorm:"size:16"`                                     // Card number for card-based payments (only for card method)
	ExpiryDate      string    `gorm:"size:5"`                                      // Expiry date for cards (e.g., "12/25", only for card method)
	AccountNumber   string    `gorm:"size:20"`                                     // Account number for bank transfer (only for bank transfer method)
	Details         string    `gorm:"size:255;not null"`                           // Details (tokenized or masked payment info)
	Status          string    `gorm:"size:20;not null"`                            // Status of the payment method (e.g., active, inactive)
	CreatedAt       time.Time `gorm:"autoCreateTime"`                              // Timestamp for when the payment method was created
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`                              // Timestamp for when the payment method was last updated
}

// TableName explicitly sets the table name to "PaymentMethods"
func (PaymentMethod) TableName() string {
	return "PaymentMethods"
}
