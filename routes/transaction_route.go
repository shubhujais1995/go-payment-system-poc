package routes

import (
	"poc/controller"
	"poc/middleware"
	"poc/services"

	"github.com/kataras/iris/v12"
)

func RegisterTransactionRoutes(app *iris.Application, svc *services.TransactionService) {
	// Protected routes for transactions
	auth := app.Party("/transactions", middleware.AuthMiddleware)
	{
		auth.Post("/", func(ctx iris.Context) {
			controller.CreateTransactionHandler(svc, ctx)
		})
		auth.Get("/", func(ctx iris.Context) {
			controller.ListTransactionsHandler(svc, ctx)
		})
	}
}
