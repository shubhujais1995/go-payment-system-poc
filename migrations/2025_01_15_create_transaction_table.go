package migrations

import (
	"fmt"
	"log"
	"poc/model"

	"gorm.io/gorm"
)

// MigrateUp - Migrate all tables up (create or update)
func MigrateUp(db *gorm.DB) error {
	// Migrate Users table
	// if err := migrateTable(db, &model.User{}, "Users"); err != nil {
	// 	return err
	// }

	// // Migrate Payers table
	// if err := migrateTable(db, &model.Payer{}, "Payers"); err != nil {
	// 	return err
	// }

	// // Migrate Payees table
	// if err := migrateTable(db, &model.Payee{}, "Payees"); err != nil {
	// 	return err
	// }

	// // Migrate Transactions table
	// if err := migrateTable(db, &model.Transaction{}, "Transactions"); err != nil {
	// 	return err
	// }

	// // Migrate Transactions table
	// if err := migrateTable(db, &model.AuditLog{}, "AuditLog"); err != nil {
	// 	return err
	// }

	// Migrate Transactions table
	if err := migrateTable(db, &model.PaymentMethod{}, "PaymentMethod"); err != nil {
		return err
	}

	// Add additional tables here as needed

	log.Println("All tables migrated successfully.")
	return nil
}

// MigrateDown - Rollback all tables (drop tables)
func MigrateDown(db *gorm.DB) error {
	// Drop Transactions table
	// if err := dropTable(db, &model.PaymentMethod{}, "PaymentMethod"); err != nil {
	// 	return err
	// }

	// // Drop Transactions table
	// if err := dropTable(db, &model.AuditLog{}, "AuditLog"); err != nil {
	// 	return err
	// }

	// Drop Transactions table
	if err := dropTable(db, &model.Transaction{}, "Transactions"); err != nil {
		return err
	}

	// // Drop Payees table
	// if err := dropTable(db, &model.Payee{}, "Payees"); err != nil {
	// 	return err
	// }

	// // Drop Payers table
	// if err := dropTable(db, &model.Payer{}, "Payers"); err != nil {
	// 	return err
	// }

	// // Drop Users table
	// if err := dropTable(db, &model.User{}, "Users"); err != nil {
	// 	return err
	// }

	// Drop additional tables as needed

	log.Println("Transactions tables dropped successfully.")
	return nil
}

// migrateTable automatically creates or updates the table schema
func migrateTable(db *gorm.DB, model interface{}, tableName string) error {
	// Perform AutoMigrate to create or update the table schema
	// log.Printf("Migrating table '%s'...\n", tableName)
	if err := db.AutoMigrate(model); err != nil {
		return fmt.Errorf("failed to migrate table '%s': %v", tableName, err)
	}

	// err = db.AutoMigrate(&model.Transaction{})
	// if err != nil {
	// 	panic("failed to migrate the database")
	// }

	// // For the Users table, create a unique index on the Email column
	// if tableName == "Users" {
	// 	// Create a unique index on the Email column
	// 	err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_Users_email ON Users (Email)`).Error
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create unique index on Users.Email: %v", err)
	// 	}
	// }

	log.Printf("Table '%s' migrated successfully.\n", tableName)
	return nil
}

// dropTable safely drops the table (and data will be lost)
func dropTable(db *gorm.DB, model interface{}, tableName string) error {
	if err := db.Migrator().DropTable(model); err != nil {
		return fmt.Errorf("failed to drop table '%s': %v", tableName, err)
	}
	log.Printf("Table '%s' dropped successfully.\n", tableName)
	return nil
}
