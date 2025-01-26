// package routes

// import (
// 	"poc/controller"
// 	"poc/middleware"
// 	"poc/services"

// 	"github.com/kataras/iris/v12"
// )

// // RegisterAuthRoutes registers the authentication-related routes.
// func RegisterAuthRoutes(app *iris.Application, svc *services.UserService) {
// 	// Create a new instance of UserController
// 	userController := &controller.UserController{
// 		UserService: svc,
// 	}

// 	// Routes
// 	app.Post("/signup", userController.Signup) // For signup
// 	app.Post("/login", userController.Login)   // For login
// 	app.Post("/logout", userController.Logout) // For logout

//		// Protected routes (example)
//		auth := app.Party("/profile", middleware.AuthMiddleware)
//		auth.Get("/", func(ctx iris.Context) {
//			userID := ctx.Values().GetString("UserID")
//			// Fetch and return the authenticated user's data
//			ctx.JSON(map[string]string{"message": "Authenticated user", "UserID": userID})
//		})
//	}
package routes

import (
	"poc/controller"
	"poc/middleware"
	"poc/services"

	"github.com/kataras/iris/v12"
)

// RegisterAuthRoutes registers the authentication-related routes.
func RegisterAuthRoutes(app *iris.Application, svc *services.UserService) {
	// Create a new instance of UserController
	userController := &controller.UserController{
		UserService: svc,
	}

	// Public routes
	app.Post("/signup", userController.Signup) // For signup
	app.Post("/login", userController.Login)   // For login
	app.Post("/logout", userController.Logout) // For logout

	// Protected routes
	auth := app.Party("/", middleware.AuthMiddleware)

	// Update user details
	auth.Put("/user", userController.UpdateUser)

	// Update payer balance
	auth.Put("/payer", userController.UpdatePayer)

	// Update payee balance
	auth.Put("/payee", userController.UpdatePayee)

	// Example protected route
	auth.Get("/profile", func(ctx iris.Context) {
		userID := ctx.Values().GetString("UserID")
		// Fetch and return the authenticated user's data
		ctx.JSON(map[string]string{"message": "Authenticated user", "UserID": userID})
	})
}
