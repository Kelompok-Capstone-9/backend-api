package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddClassRoutes(e *echo.Echo) {
	// for users
	e.GET("/classes", controllers.GetClassesController)
	e.GET("/classes/:id", controllers.GetClassByIDController)

	// for administrators
	classJWT := e.Group("/admin/classes")
	classJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	classJWT.GET("", controllers.GetClassesController)
	classJWT.GET("/:id", controllers.GetClassByIDController)
	classJWT.POST("", controllers.CreateClassController)
	classJWT.PUT("/:id", controllers.EditClassController)
	classJWT.DELETE("/:id", controllers.DeleteClassController)
	classJWT.POST("/banner/:id", controllers.UploadClassImageController) // for upload class image banner
}
