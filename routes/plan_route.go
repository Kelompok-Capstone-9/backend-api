package routes

import (
	"gofit-api/controllers"

	"github.com/labstack/echo/v4"
)

func AddPlanRoutes(e *echo.Echo) {

	e.GET("/plans", controllers.GetPlansController)
	e.GET("/plans/:id", controllers.GetPlanController)
	e.POST("/plans", controllers.CreatePlanController)
	e.PUT("/plans/:id", controllers.UpdatePlanController)
	e.DELETE("/plans/:id", controllers.DeletePlanController)
}
