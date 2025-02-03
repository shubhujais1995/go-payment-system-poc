package initializer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock implementation of a GORM database connection
type MockDB struct {
	mock.Mock
	*gorm.DB
}

// working utc
// Test InitializeGORMSpannerClient with a successful connection
func TestInitializeGORMSpannerClient_Success(t *testing.T) {
	// Reset dbInstance before running the test
	dbInstance = nil

	// Mock environment variables
	os.Setenv("DB_PROJECT_ID", "mock-project")
	os.Setenv("DB_INSTANCE_ID", "mock-instance")
	os.Setenv("DB_NAME", "mock-db")

	// Mock the database connection instead of opening a real connection
	mockDB := &gorm.DB{} // Create a fake GORM DB object
	dbInstance = mockDB  // Assign the mock DB instance

	// Call the function
	db, err := InitializeGORMSpannerClient()

	// Assertions
	assert.Nil(t, err, "Expected no error during database initialization")
	assert.NotNil(t, db, "Expected a valid DB instance")
	assert.Equal(t, dbInstance, db, "dbInstance should be set and match the returned DB instance")
}
