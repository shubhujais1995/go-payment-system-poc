package model

import "time"

// User represents the basic account and authentication information for a user.
type User struct {
	UserID       string    `gorm:"primaryKey;column:UserID"` // Unique ID for the user
	Email        string    `gorm:"unique;column:Email"`      // Unique email for user authentication
	PasswordHash string    `gorm:"column:PasswordHash"`      // Hashed password for authentication
	FirstName    string    `gorm:"column:FirstName"`         // User's first name
	LastName     string    `gorm:"column:LastName"`          // User's last name
	IsVerified   bool      `gorm:"column:IsVerified"`        // Indicates if the user is verified
	CreatedAt    time.Time `gorm:"column:CreatedAt"`         // When the user was created
	UpdatedAt    time.Time `gorm:"column:UpdatedAt"`         // Last update timestamp for the user
}

// TableName explicitly sets the table name to "Users" (case-sensitive).
func (User) TableName() string {
	return "Users" // Use the exact case-sensitive name you want
}

type EntityWithBalance struct {
	ID string

	Balance float64
}
