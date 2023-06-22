package routes

import (
	"gofit-api/configs"
	"gofit-api/middlewares"
	"gofit-api/models"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	jwtConfig echojwt.Config
)

func LoadJwtConfig() {
	jwtConfig = echojwt.Config{
		SigningKey: []byte(configs.AppConfig.JWTKey),
		ErrorHandler: func(c echo.Context, err error) error {
			if err != nil {
				var response models.GeneralResponse
				response.StatusCode = http.StatusBadRequest
				response.Message = err.Error()
				response.ErrorReason = "no token been inputed"
				return c.JSON(response.StatusCode, response)
			}
			if ve, ok := err.(*jwt.ValidationError); ok {
				var response models.GeneralResponse
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					response.StatusCode = http.StatusBadRequest
					response.Message = jwt.ErrTokenMalformed.Error()
					response.ErrorReason = "token is malformed"
					return c.JSON(response.StatusCode, response)
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					response.StatusCode = http.StatusUnauthorized
					response.Message = jwt.ErrTokenExpired.Error()
					response.ErrorReason = "token is expired"
					return c.JSON(response.StatusCode, response)
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					response.StatusCode = http.StatusUnauthorized
					response.Message = jwt.ErrTokenNotValidYet.Error()
					response.ErrorReason = "token is not valid yet"
					return c.JSON(response.StatusCode, response)
				} else {
					return c.JSON(http.StatusBadRequest, "token is invalid")
				}
			}
			return nil
		},
	}
}

func InitRoute(e *echo.Echo) {
	middlewares.Logger(e)
	e.Use(middleware.CORS())
	e.Static("/assets", "assets")
	LoadJwtConfig()
	AddUserRoutes(e)
	AddPlanRoutes(e)
	AddMembershipRoutes(e)
	AddLocationRoutes(e)
	AddClassRoutes(e)
	AddClassPackageRoutes(e)
	AddClassTicketRoutes(e)
	AddHealthtipRoutes(e)
}
