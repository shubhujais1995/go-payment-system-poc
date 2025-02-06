package initializer

import (
	"fmt"
	"log"

	spannergorm "github.com/googleapis/go-gorm-spanner"
	_ "github.com/googleapis/go-sql-spanner"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbOpen = gorm.Open
var dbInstance *gorm.DB

func InitializeGORMSpannerClient() (*gorm.DB, error) {
	if dbInstance != nil {
		log.Println("Database already initialized, returning existing instance.")
		return dbInstance, nil
	}

	projectID := GetEnv("DB_PROJECT_ID")
	instanceID := GetEnv("DB_INSTANCE_ID")
	dbName := GetEnv("DB_NAME")

	dsn := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceID, dbName)

	db, err := dbOpen(spannergorm.New(spannergorm.Config{
		DriverName: "spanner",
		DSN:        dsn,
	}), &gorm.Config{
		PrepareStmt:                      true,
		IgnoreRelationshipsWhenMigrating: true,
		Logger:                           logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Printf("Failed to open database: %v", err)
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	dbInstance = db
	log.Printf("Successfully connected to Spanner database: %s", dsn)
	return db, nil
}

func GetDB() *gorm.DB {
	if dbInstance == nil {
		log.Printf("Database not initialized. Call InitializeGORMSpannerClient first.")
		panic("Database not initialized")
	}
	return dbInstance
}

func SetDBForTest(db *gorm.DB) {
	dbInstance = db
}

// // runMigrateUp handles the logic for running migrations to create or update the tables
// func runMigrateUp(db *gorm.DB) {
// 	if err := migrations.MigrateUp(db); err != nil {
// 		log.Fatalf("Migration Up failed: %v", err)
// 	}
// 	fmt.Println("Migration Up completed successfully.")
// }

// // runMigrateDown handles the logic for rolling back the migrations (dropping tables)
// func runMigrateDown(db *gorm.DB) {
// 	if err := migrations.MigrateDown(db); err != nil {
// 		log.Fatalf("Migration Down failed: %v", err)
// 	}
// 	fmt.Println("Migration Down completed successfully.")
// }
