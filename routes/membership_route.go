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
	membershipJWT.POST("/join/:plan_id", controllers.JoinMembershipController)

	adminMembershipJWT := e.Group("/admin/memberships")
	adminMembershipJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	adminMembershipJWT.GET("", controllers.GetMembershipsController)
	adminMembershipJWT.GET("/:id", controllers.GetMembershipController)
	adminMembershipJWT.POST("", controllers.CreateMembershipController)
	adminMembershipJWT.PUT("/:id", controllers.UpdateMembershipController)
	adminMembershipJWT.DELETE("/:id", controllers.DeleteMembershipController)
}
