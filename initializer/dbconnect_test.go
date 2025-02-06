package initializer

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInitializeGORMSpannerClient(t *testing.T) {
	// Reset dbInstance before each test
	dbInstance = nil

	// Store original dbOpen function and environment variables
	originalDBOpen := dbOpen
	originalProjectID := os.Getenv("DB_PROJECT_ID")
	originalInstanceID := os.Getenv("DB_INSTANCE_ID")
	originalDBName := os.Getenv("DB_NAME")

	// Restore original values after tests
	defer func() {
		dbOpen = originalDBOpen
		os.Setenv("DB_PROJECT_ID", originalProjectID)
		os.Setenv("DB_INSTANCE_ID", originalInstanceID)
		os.Setenv("DB_NAME", originalDBName)
	}()

	t.Run("successful initialization", func(t *testing.T) {
		dbInstance = nil
		os.Setenv("DB_PROJECT_ID", "test-project")
		os.Setenv("DB_INSTANCE_ID", "test-instance")
		os.Setenv("DB_NAME", "test-db")

		dbOpen = func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error) {
			return &gorm.DB{}, nil
		}

		db, err := InitializeGORMSpannerClient()
		assert.NoError(t, err)
		assert.NotNil(t, db)
		assert.Equal(t, dbInstance, db)
	})

	t.Run("return existing instance if already initialized", func(t *testing.T) {
		mockInstance := &gorm.DB{}
		dbInstance = mockInstance

		db, err := InitializeGORMSpannerClient()
		assert.NoError(t, err)
		assert.Equal(t, mockInstance, db)
	})

	t.Run("handle database connection error", func(t *testing.T) {
		dbInstance = nil
		os.Setenv("DB_PROJECT_ID", "test-project")
		os.Setenv("DB_INSTANCE_ID", "test-instance")
		os.Setenv("DB_NAME", "test-db")

		dbOpen = func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error) {
			return nil, errors.New("connection failed")
		}

		db, err := InitializeGORMSpannerClient()
		assert.Error(t, err)
		assert.Nil(t, db)
		assert.Contains(t, err.Error(), "failed to open database")
	})
}

func TestGetDB(t *testing.T) {
	t.Run("get initialized database", func(t *testing.T) {
		mockInstance := &gorm.DB{}
		dbInstance = mockInstance

		db := GetDB()
		assert.Equal(t, mockInstance, db)
	})

	t.Run("panic when database not initialized", func(t *testing.T) {
		dbInstance = nil

		assert.Panics(t, func() {
			GetDB()
		})
	})
}

func TestSetDBForTest(t *testing.T) {
	dbInstance = nil
	mockDB := &gorm.DB{}
	SetDBForTest(mockDB)
	assert.Equal(t, mockDB, dbInstance)
}
