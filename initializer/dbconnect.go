package initializer

import (
	"fmt"
	"log"
	"poc/migrations"

	spannergorm "github.com/googleapis/go-gorm-spanner"
	_ "github.com/googleapis/go-sql-spanner" // Import the Spanner SQL driver
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// dbOpen is a wrapper for gorm.Open, which can be mocked in tests
var dbOpen = gorm.Open

var dbInstance *gorm.DB

// InitializeGORMSpannerClient initializes GORM with Google Cloud Spanner.
func InitializeGORMSpannerClient() (*gorm.DB, error) {
	// If the DB instance is already initialized, return it
	if dbInstance != nil {
		log.Println("Database already initialized, returning existing instance.")
		return dbInstance, nil
	}

	// Get the database connection details from the environment
	projectID := GetEnv("DB_PROJECT_ID")
	instanceID := GetEnv("DB_INSTANCE_ID")
	dbName := GetEnv("DB_NAME")

	// Build the Spanner connection string
	dsn := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, dbName)

	// Connect to Spanner using the GORM driver for Spanner
	db, err := dbOpen(spannergorm.New(spannergorm.Config{
		DriverName: "spanner",
		DSN:        dsn,
	}), &gorm.Config{
		PrepareStmt:                      true,
		IgnoreRelationshipsWhenMigrating: true,
		Logger:                           logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Assign the new connection to the global instance
	dbInstance = db

	// Log success
	log.Printf("Successfully connected to Spanner database: %s", dsn)
	return db, nil
}

// package initializer

// import (
// 	"fmt"
// 	"log"
// 	"poc/migrations"

// 	spannergorm "github.com/googleapis/go-gorm-spanner"
// 	_ "github.com/googleapis/go-sql-spanner" // Import the Spanner SQL driver
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/logger"
// )

// var dbInstance *gorm.DB

// // InitializeGORMSpannerClient initializes GORM with Google Cloud Spanner.
// func InitializeGORMSpannerClient() (*gorm.DB, error) {
// 	// If the DB instance is already initialized, return it
// 	if dbInstance != nil {
// 		log.Println("Database already initialized, returning existing instance.")
// 		return dbInstance, nil
// 	}

// 	// Get the database connection details from the environment
// 	projectID := GetEnv("DB_PROJECT_ID")
// 	instanceID := GetEnv("DB_INSTANCE_ID")
// 	dbName := GetEnv("DB_NAME")
// 	fmt.Println("run test => Project id ", projectID, " inst id ", instanceID, " dbName ", dbName)

// 	// Build the Spanner connection string
// 	dsn := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, dbName)

// 	fmt.Println(dsn, " = dsn")
// 	// Connect to Spanner using the GORM driver for Spanner
// 	db, err := gorm.Open(spannergorm.New(spannergorm.Config{
// 		DriverName: "spanner",
// 		DSN:        dsn,
// 	}), &gorm.Config{
// 		PrepareStmt:                      true,
// 		IgnoreRelationshipsWhenMigrating: true,
// 		Logger:                           logger.Default.LogMode(logger.Error),
// 	})

// 	if err != nil {
// 		log.Fatalf("Failed to open database: %v", err)
// 		return nil, fmt.Errorf("failed to open database: %w", err)
// 	}

// 	// Assign the new connection to the global instance
// 	dbInstance = db

// 	// Log success
// 	log.Printf("Successfully connected to Spanner database: %s", dsn)
// 	return db, nil
// }

func GetDB() *gorm.DB {
	if dbInstance == nil {
		log.Fatal("Database not initialized. Call InitializeGORMSpannerClient first.")
	}
	return dbInstance
}

// runMigrateUp handles the logic for running migrations to create or update the tables
func runMigrateUp(db *gorm.DB) {
	if err := migrations.MigrateUp(db); err != nil {
		log.Fatalf("Migration Up failed: %v", err)
	}
	fmt.Println("Migration Up completed successfully.")
}

// runMigrateDown handles the logic for rolling back the migrations (dropping tables)
func runMigrateDown(db *gorm.DB) {
	if err := migrations.MigrateDown(db); err != nil {
		log.Fatalf("Migration Down failed: %v", err)
	}
	fmt.Println("Migration Down completed successfully.")
}
