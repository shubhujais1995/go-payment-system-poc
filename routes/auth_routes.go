package routes

import (
	"poc/controller"
	"poc/middleware"

	"github.com/kataras/iris/v12"
)

// RegisterAuthRoutes registers the authentication routes.
func RegisterAuthRoutes(app *iris.Application) {
    // Public routes
    app.Post("/signup", controller.Signup) // For signup
    app.Post("/login", controller.Login)   // For login
    app.Post("/logout", controller.Logout) // For logout

    // Protected routes
    auth := app.Party("/", middleware.AuthMiddleware)

    // Update user details
    auth.Put("/user", controller.UpdateUser)

    // Update payer balance
    auth.Put("/payer", controller.UpdatePayer)

    // Update payee balance
    auth.Put("/payee", controller.UpdatePayee)

    // Example protected route
    auth.Get("/profile", func(ctx iris.Context) {
        userID := ctx.Values().GetString("UserID")
        // Fetch and return the authenticated user's data
        ctx.JSON(map[string]string{"message": "Authenticated user", "UserID": userID})
    })
}