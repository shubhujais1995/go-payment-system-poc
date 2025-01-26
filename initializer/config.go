package initializer

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig loads environment variables from a .env file
func LoadConfig() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Ensure the required environment variables are set
	if os.Getenv("DB_PROJECT_ID") == "" {
		log.Fatal("Missing DB_PROJECT_ID environment variable")
	}
	if os.Getenv("DB_INSTANCE_ID") == "" {
		log.Fatal("Missing DB_INSTANCE_ID environment variable")
	}
	if os.Getenv("DB_NAME") == "" {
		log.Fatal("Missing DB_NAME environment variable")
	}

	// Ensure GOOGLE_APPLICATION_CREDENTIALS is set if running locally
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		log.Fatal("Missing GOOGLE_APPLICATION_CREDENTIALS environment variable")
	}
}

// GetEnv is a helper function to fetch environment variables
func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}
