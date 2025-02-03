package routes

import (
	"poc/controller"
	"poc/middleware"

	"github.com/kataras/iris/v12"
)

func RegisterTransactionRoutes(app *iris.Application) {
	// Protected routes for transactions
	auth := app.Party("/transactions", middleware.AuthMiddleware)
	{
		auth.Post("/", func(ctx iris.Context) {
			controller.CreateTransactionHandler(ctx)
		})
		auth.Get("/", func(ctx iris.Context) {
			controller.ListTransactionsHandler(ctx)
		})
	}
}
