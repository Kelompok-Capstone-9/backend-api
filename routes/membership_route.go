package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddMembershipRoutes(e *echo.Echo) {
	membershipJWT := e.Group("/memberships")
	membershipJWT.Use(echojwt.WithConfig(jwtConfig))
	membershipJWT.Use(m.IsAdmin)

	membershipJWT.GET("", controllers.GetMembershipsController)
	membershipJWT.GET("/:id", controllers.GetMembershipController, m.IsSameUser)
	membershipJWT.POST("", controllers.CreateMembershipController)
	membershipJWT.PUT("/:id", controllers.UpdateMembershipController)
	membershipJWT.DELETE("/:id", controllers.DeleteMembershipController)
}
