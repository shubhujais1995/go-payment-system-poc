package main

import (
	"fmt"
	"log"
	"os"

	"poc/initializer"
	"poc/routes"

	"github.com/kataras/iris/v12"
)

func main() {
	// Load configuration from .env file
	initializer.LoadConfig()

	// Initialize the GORM client (Spanner)
	_, err := initializer.InitializeGORMSpannerClient()
	if err != nil {
		log.Fatalf("Failed to initialize GORM Spanner client: %v", err)
	}

	// Create an Iris application instance
	app := iris.New()

	app.HandleDir("/", iris.Dir("."))

	// Register routes for user, transaction, and payment method
	routes.RegisterAuthRoutes(app)
	routes.RegisterPaymentRoutes(app) // Add this to register payment method routes
	routes.RegisterTransactionRoutes(app)

	// Define the server port (default to 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	fmt.Printf("Server is running on port %s...\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
