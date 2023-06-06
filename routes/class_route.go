package routes

import (
	"gofit-api/controllers"

	"github.com/labstack/echo/v4"
)

func AddClassRoutes(e *echo.Echo) {

	e.GET("/classes/all", controllers.GetClassesController)
	// e.GET("/classes/:id", controllers.GetInstuctorController)

	// instructorJWT := e.Group("/instructors")
	// instructorJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	// instructorJWT.POST("", controllers.CreateInstructorController)
	// instructorJWT.PUT("/:id", controllers.EditInstructorController)
	// instructorJWT.DELETE("/:id", controllers.DeleteInstructorController)
}
