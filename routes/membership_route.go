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
	membershipJWT.GET("/mymembership", controllers.MyMembershipController)

	membershipJWT.GET("", controllers.GetMembershipsController, m.IsAdmin)
	membershipJWT.GET("/:id", controllers.GetMembershipController, m.IsAdmin)
	membershipJWT.POST("", controllers.CreateMembershipController, m.IsAdmin)
	membershipJWT.PUT("/:id", controllers.UpdateMembershipController, m.IsAdmin)
	membershipJWT.DELETE("/:id", controllers.DeleteMembershipController, m.IsAdmin)

	membershipJWT.POST("/join/:plan_id", controllers.JoinMembershipController)

}
