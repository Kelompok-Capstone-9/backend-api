package routes

import (
	"gofit-api/controllers"
	m "gofit-api/middlewares"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func AddUserRoutes(e *echo.Echo) {

	e.POST("/register", controllers.CreateUserController)
	e.POST("/login", controllers.LoginUserController)
	e.POST("/forgot", controllers.ForgotPasswordController)

	userJWT := e.Group("/users")
	userJWT.Use(echojwt.WithConfig(jwtConfig))
	userJWT.GET("", controllers.GetUsersController, m.IsAdmin) // for admin
	userJWT.GET("/:id", controllers.GetUserController, m.IsSameUser)
	userJWT.PUT("/:id", controllers.EditUserController, m.IsSameUser)
	userJWT.POST("/profile/:id", controllers.UploadProfilePictureController, m.IsSameUser) // for upload profile picture
	userJWT.DELETE("/:id", controllers.DeleteUserController, m.IsAdmin)                    // for admin
}
