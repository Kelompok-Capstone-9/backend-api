package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddClassPackageRoutes(e *echo.Echo) {

	// for administrator
	classPackageJWT := e.Group("/admin/classes/packages")
	classPackageJWT.Use(echojwt.WithConfig(jwtConfig), m.IsAdmin)
	classPackageJWT.GET("", controllers.GetClassPackagesController) // with params
	classPackageJWT.GET("/:id", controllers.GetClassPackageByIDController)
	classPackageJWT.POST("", controllers.CreateClassPackageController)
	classPackageJWT.PUT("/:id", controllers.EditClassPackageController)
	classPackageJWT.DELETE("/:id", controllers.DeleteClassPackageController)
}
