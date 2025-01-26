package model

import "time"

// Payee represents a user who can receive money (payee details).
type Payee struct {
	PayeeID   string    `gorm:"primaryKey;column:PayeeID"` // Unique payee ID
	UserID    string    `gorm:"not null;index"`            // Foreign key to the User table
	Name      string    `gorm:"column:Name"`               // Name of the payee (individual or business)
	Email     string    `gorm:"column:Email"`              // Contact email for the payee
	Address   string    `gorm:"column:Address"`            // Physical address (optional)
	Balance   float64   `gorm:"column:Balance"`            // Available balance for the payee
	Status    string    `gorm:"column:Status"`             // Account status (active, inactive, suspended)
	CreatedAt time.Time `gorm:"column:CreatedAt"`          // Timestamp for when the payee record was created
	UpdatedAt time.Time `gorm:"column:UpdatedAt"`          // Timestamp for when the payee record was last updated
}

// TableName explicitly sets the table name to "Payees".
func (Payee) TableName() string {
	return "Payees"
}
