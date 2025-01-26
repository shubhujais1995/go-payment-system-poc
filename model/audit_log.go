package model

import "time"

// AuditLog represents an entry in the audit log, tracking actions on transactions.
type AuditLog struct {
	AuditLogID    string    `gorm:"primaryKey;size:36"` // Unique identifier for the audit log
	TransactionID string    `gorm:"size:36;index"`      // Associated transaction ID
	Action        string    `gorm:"size:255;not null"`  // Description of the action performed
	CreatedAt     time.Time `gorm:"autoCreateTime"`     // Timestamp when the log entry was created
	Details       string    `gorm:"size:255"`           // Additional details of the action
}

// TableName explicitly sets the table name to "AuditLogs"
func (AuditLog) TableName() string {
	return "AuditLogs"
}
