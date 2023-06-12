package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddClassPackageRoutes(e *echo.Echo) {

	e.GET("/classes/packages/all", controllers.GetClassPackgesController)
	// e.GET("/classes/:id", controllers.GetClassByIDController)

	// for administrator
	classJWT := e.Group("/classes/package")
	classJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	classJWT.POST("", controllers.CreateClassController)
	classJWT.PUT("/:id", controllers.EditClassController)
	classJWT.DELETE("/:id", controllers.DeleteClassController)
}
