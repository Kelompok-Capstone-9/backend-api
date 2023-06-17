package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddLocationRoutes(e *echo.Echo) {
	e.GET("/locations", controllers.GetLocationsController)
	e.GET("/locations/:id", controllers.GetLocationByIDController)

	// for administrator
	locationJWT := e.Group("/admin/locations")
	locationJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	locationJWT.POST("", controllers.CreateLocationController)
	locationJWT.PUT("/:id", controllers.EditLocationController)
	locationJWT.DELETE("/:id", controllers.DeleteLocationController)
}
