package routes

import (
	"gofit-api/controllers"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddHealthtipRoutes(e *echo.Echo) {

	healthtipJWT := e.Group("/healthtips")
	healthtipJWT.Use(echojwt.WithConfig(jwtConfig))
	healthtipJWT.GET("", controllers.GetHealthtipsController)
	healthtipJWT.GET("/:id", controllers.GetHealthtipByIDController)
	healthtipJWT.POST("", controllers.CreateHealthtipController)
	healthtipJWT.PUT("/:id", controllers.EditHealthtipController)
	healthtipJWT.DELETE("/:id", controllers.DeleteHealthtipController)

}
