package middlewares

import (
	"errors"
	"fmt"
	"gofit-api/configs"
	"gofit-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// (userID int, email string, isAdmin bool) => for admin role
func CreateToken(userID int, email string, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{}
	claims["userID"] = userID
	claims["email"] = email
	claims["isAdmin"] = isAdmin
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()
	fmt.Println(claims["exp"])

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.AppConfig.JWTKey))
}

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isAdmin := claims["isAdmin"].(bool)

		if !isAdmin {
			var response models.GeneralResponse
			response.StatusCode = http.StatusUnauthorized
			response.Message = "Unauthorized"
			response.ErrorReason = "cannot access without permission"
			return c.JSON(response.StatusCode, response)
		}

		return next(c)
	}
}

func IsSameUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isAdmin := claims["isAdmin"].(bool)

		if isAdmin {
			return next(c)
		}

		tokenUserID := claims["userID"].(float64)
		paramUserID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			var response models.GeneralResponse
			response.StatusCode = http.StatusBadRequest
			response.Message = err.Error()
			response.ErrorReason = "invalid id param"
			return c.JSON(response.StatusCode, response)
		}

		isSameUser := int(tokenUserID) == paramUserID
		if !isSameUser {
			var response models.GeneralResponse
			response.StatusCode = http.StatusUnauthorized
			response.Message = "unauthorized"
			response.ErrorReason = "not the same user as id inputed in parameter"
			return c.JSON(response.StatusCode, response)
		}
		return next(c)
	}
}

func ExtractTokenInfo(e echo.Context) (models.TokenInfo, error) {
	var tokenInfo models.TokenInfo
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		tokenInfo.UserID = int(claims["userID"].(float64))
		tokenInfo.Email = claims["email"].(string)
		tokenInfo.IsAdmin = claims["isAdmin"].(bool)
		tokenInfo.Expired = claims["exp"].(time.Time)
		return tokenInfo, nil
	}
	return tokenInfo, errors.New("no token found")
}

// func ExtractTokenUserID(e echo.Context) float64 {
// 	user := e.Get("user").(*jwt.Token)
// 	if user.Valid {
// 		claims := user.Claims.(jwt.MapClaims)
// 		userId := claims["userID"].(float64)
// 		return userId
// 	}
// 	return 0
// }

// func ExtractTokenIsAdmin(e echo.Context) bool {
// 	user := e.Get("user").(*jwt.Token)
// 	if user.Valid {
// 		claims := user.Claims.(jwt.MapClaims)
// 		isAdmin := claims["isAdmin"].(bool)
// 		return isAdmin
// 	}
// 	return false
// }
