package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddLocationRoutes(e *echo.Echo) {
	locationJWT := e.Group("/locations")
	locationJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	locationJWT.GET("/all", controllers.GetLocationsController)
	locationJWT.GET("/:id", controllers.GetLocationByIDController)
	locationJWT.POST("", controllers.CreateLocationController)
	locationJWT.PUT("/:id", controllers.EditLocationController)
	locationJWT.DELETE("/:id", controllers.DeleteLocationController)
}
