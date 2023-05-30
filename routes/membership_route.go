package routes

import (
	"gofit-api/controllers"

	"github.com/labstack/echo/v4"
)

func AddMembershipRoutes(e *echo.Echo) {

	e.GET("/memberships", controllers.GetMembershipsController)
	e.GET("/memberships/:id", controllers.GetMembershipController)
	e.POST("/memberships", controllers.CreateMembershipController)
	e.PUT("/memberships/:id", controllers.UpdateMembershipController)
	e.DELETE("/memberships/:id", controllers.DeleteMembershipController)
}
