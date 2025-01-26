package model

import "time"

// Payer represents a user who can send money (payer details).
type Payer struct {
	PayerID         string    `gorm:"primaryKey;column:PayerID"` // Unique payer ID
	UserID          string    `gorm:"not null;index"`            // Foreign key to the User table
	Name            string    `gorm:"column:Name"`               // Name of the payer (individual or business)
	Email           string    `gorm:"column:Email"`              // Contact email for the payer
	PhoneNumber     string    `gorm:"column:PhoneNumber"`        // Phone number (optional)
	Address         string    `gorm:"column:Address"`            // Physical address (optional)
	PaymentMethodID string    `gorm:"column:PaymentMethodID"`    // Payment method ID (e.g., card, bank account, wallet)
	Balance         float64   `gorm:"column:Balance"`            // Available balance for the payer
	Status          string    `gorm:"column:Status"`             // Account status (active, inactive, suspended)
	CreatedAt       time.Time `gorm:"column:CreatedAt"`          // Timestamp for when the payer record was created
	UpdatedAt       time.Time `gorm:"column:UpdatedAt"`          // Timestamp for when the payer record was last updated
}

// TableName explicitly sets the table name to "Payers".
func (Payer) TableName() string {
	return "Payers"
}
