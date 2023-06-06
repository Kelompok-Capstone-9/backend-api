package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddPlanRoutes(e *echo.Echo) {
	e.GET("/plans/all", controllers.GetPlansController)
	e.GET("/plans/:id", controllers.GetPlanController)

	planJWT := e.Group("/plans")
	planJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	planJWT.POST("", controllers.CreatePlanController)
	planJWT.PUT("/:id", controllers.UpdatePlanController)
	planJWT.DELETE("/:id", controllers.DeletePlanController)
}
