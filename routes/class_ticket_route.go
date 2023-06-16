package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddClassTicketRoutes(e *echo.Echo) {
	classTicketJWT := e.Group("/classes/tickets")
	classTicketJWT.Use(echojwt.WithConfig(jwtConfig))

	// for users
	classTicketJWT.GET("/mytickets", controllers.GetMyTicketsController)              // my tickets (get tickets by user id)
	classTicketJWT.GET("/mytickets/:id", controllers.GetMyTicketDetailController)     // get my ticket details
	classTicketJWT.GET("/mytickets/cancel/:id", controllers.CancelMyTicketController) // cancel my ticket/booking

	// for administrator
	ticketAdminJWT := e.Group("/admin/classes/tickets")
	ticketAdminJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	ticketAdminJWT.GET("/:id", controllers.GetClassTicketByIDController)
	ticketAdminJWT.GET("", controllers.GetClassTicketsController) // with params
	ticketAdminJWT.POST("", controllers.CreateClassTicketController)
	ticketAdminJWT.PUT("/:id", controllers.EditClassTicketController)
	ticketAdminJWT.DELETE("/:id", controllers.DeleteClassTicketController)
}
