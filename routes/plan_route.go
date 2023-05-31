package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddPlanRoutes(e *echo.Echo) {
	planJWT := e.Group("/plans")
	planJWT.Use(echojwt.WithConfig(jwtConfig))
	planJWT.Use(m.IsAdmin)

	planJWT.GET("", controllers.GetPlansController)
	planJWT.GET("/:id", controllers.GetPlanController)
	planJWT.POST("", controllers.CreatePlanController)
	planJWT.PUT("/:id", controllers.UpdatePlanController)
	planJWT.DELETE("/:id", controllers.DeletePlanController)
}
