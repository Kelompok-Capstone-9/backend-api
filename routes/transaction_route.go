package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddTransactionRoutes(e *echo.Echo) {

	e.POST("/transactions/pay/:transaction_code", controllers.PayController)
	// transactionJWT := e.Group("/transaction")
	// transactionJWT.Use(echojwt.WithConfig(jwtConfig))

	adminTransactionJWT := e.Group("/admin/transactions")
	adminTransactionJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	adminTransactionJWT.GET("", controllers.GetTransactionsController)
	adminTransactionJWT.GET("/:id", controllers.GetTransactionController)
	adminTransactionJWT.POST("", controllers.CreateTransactionController)
	adminTransactionJWT.PUT("/:id", controllers.UpdateTransactionController)
	adminTransactionJWT.DELETE("/:id", controllers.DeleteTransactionController)
}
