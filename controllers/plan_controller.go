package controllers

import (
	"gofit-api/lib/database"
	"gofit-api/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetPlansController retrieves all plans
func GetPlansController(c echo.Context) error {
	var response models.GeneralListResponse
	var page models.Pages
	var err models.CustomError

	page.PageString = c.QueryParam("page")
	page.ConvertPageToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	offset, limit := page.CalcOffsetLimit()
	plans, totalData := database.GetPlans(offset, limit, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success("Successfully retrieved plans", page.Page, totalData, plans)
	return c.JSON(http.StatusOK, response)
}

// Get plan by ID
func GetPlanController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readablePlan models.ReadablePlan
	var planObject models.Plan

	planID := c.Param("id")

	planObject.InsertID(planID, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetPlan(&planObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readablePlan.ToPlanObject(&planObject)

	response.Success(http.StatusOK, "Successfully retrieved plan", readablePlan)
	return c.JSON(http.StatusOK, response)
}

// Create new plan
func CreatePlanController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readablePlan models.ReadablePlan
	var planObject models.Plan

	err.ErrorMessage = c.Bind(&readablePlan)
	if err.IsError() {
		err.StatusCode = http.StatusBadRequest
		err.ErrorReason = "Invalid request body"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readablePlan.ToPlanObject(&planObject)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.CreatePlan(&planObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readablePlan.ToPlanObject(&planObject)

	response.Success(http.StatusCreated, "Successfully created a new plan", readablePlan)
	return c.JSON(response.StatusCode, response)
}
