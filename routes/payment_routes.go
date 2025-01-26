package routes

import (
	"poc/controller"
	"poc/middleware"
	"poc/services"

	"github.com/kataras/iris/v12"
)

func RegisterPaymentRoutes(app *iris.Application, svc *services.PaymentMethodService) {
	// Protected routes for payment methods
	auth := app.Party("/payment-methods", middleware.AuthMiddleware) // Apply authentication middleware
	{
		// Route for creating a payment method
		auth.Post("/", func(ctx iris.Context) {
			controller.CreatePaymentMethodHandler(svc, ctx)
		})

		// Route for fetching payment method
		auth.Get("/", func(ctx iris.Context) {
			controller.GetPaymentMethodHandler(svc, ctx)
		})

		// Route for updating payment method
		auth.Put("/{paymentMethodID}", func(ctx iris.Context) {
			controller.UpdatePaymentMethodHandler(svc, ctx)
		})

		// Route for validating payment method
		auth.Post("/validate/{paymentMethodID}", func(ctx iris.Context) {
			controller.ValidatePaymentMethodHandler(svc, ctx)
		})
	}
}
