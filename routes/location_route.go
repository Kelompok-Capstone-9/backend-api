package routes

import (
	"gofit-api/controllers"

	"github.com/labstack/echo/v4"
)

func AddLocationRoutes(e *echo.Echo) {

	e.GET("/locations/all", controllers.GetLocationsController)
	// e.GET("/classes/:id", controllers.GetInstuctorController)

	// instructorJWT := e.Group("/instructors")
	// instructorJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	// instructorJWT.POST("", controllers.CreateInstructorController)
	// instructorJWT.PUT("/:id", controllers.EditInstructorController)
	// instructorJWT.DELETE("/:id", controllers.DeleteInstructorController)
}
