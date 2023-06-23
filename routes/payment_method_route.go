package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddPaymentMethodRoutes(e *echo.Echo) {
	// for users or guests
	e.GET("/transactions/payment_methods", controllers.GetPaymentMethodsController)
	e.GET("/transactions/payment_methods/:id", controllers.GetPaymentMethodController)

	// for admin
	paymentMethodJWT := e.Group("/admin/transactions/payment_methods")
	paymentMethodJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	paymentMethodJWT.POST("", controllers.CreatePaymentMethodController)
	paymentMethodJWT.PUT("/:id", controllers.UpdatePaymentMethodController)
	paymentMethodJWT.DELETE("/:id", controllers.DeletePaymentMethodController)
}
