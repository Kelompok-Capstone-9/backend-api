package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddClassTicketRoutes(e *echo.Echo) {

	// my tickets (get tickets by user id)

	e.GET("/classes/tickets/:id", controllers.GetClassTicketByIDController)

	// for administrator
	classTicketJWT := e.Group("/classes/tickets")
	classTicketJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	classTicketJWT.GET("/all", controllers.GetClassTicketsController) // with params
	classTicketJWT.POST("", controllers.CreateClassTicketController)
	classTicketJWT.PUT("/:id", controllers.EditClassTicketController)
	classTicketJWT.DELETE("/:id", controllers.DeleteClassTicketController)
}
