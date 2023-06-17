package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddPlanRoutes(e *echo.Echo) {
	// for users or guests
	e.GET("/plans/all", controllers.GetPlansController)
	e.GET("/plans/:id", controllers.GetPlanController)

	// for admin
	planJWT := e.Group("/admin/plans")
	planJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	planJWT.POST("", controllers.CreatePlanController)
	planJWT.PUT("/:id", controllers.UpdatePlanController)
	planJWT.DELETE("/:id", controllers.DeletePlanController)
}
