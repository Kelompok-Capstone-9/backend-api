package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddClassRoutes(e *echo.Echo) {

	e.GET("/classes/all", controllers.GetClassesController)
	e.GET("/classes/:id", controllers.GetClassByIDController)

	classJWT := e.Group("/classes")
	classJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	classJWT.POST("", controllers.CreateClassController)
	classJWT.PUT("/:id", controllers.EditClassController)
	classJWT.DELETE("/:id", controllers.DeleteClassController)
}
