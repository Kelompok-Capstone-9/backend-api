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
	page.ConvertPageStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	page.CalcOffsetLimit()
	plans, totalData := database.GetPlans(page.Offset, page.Limit, &err)
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

	// Extract the plan ID from the URL parameter
	planID := c.Param("id")

	// Convert the plan ID string to int
	planObject.InsertID(planID, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// Retrieve the plan from the database
	database.GetPlan(&planObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// Convert the plan object to a readable format
	readablePlan.ToReadablePlan(&planObject)

	response.Success(http.StatusOK, "Successfully retrieved plan", readablePlan)
	return c.JSON(http.StatusOK, response)
}

// Create new plan
func CreatePlanController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var planObject models.Plan

	err.ErrorMessage = c.Bind(&planObject)
	if err.IsError() {
		err.StatusCode = http.StatusBadRequest
		err.ErrorReason = "Invalid request body"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.CreatePlan(&planObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success(http.StatusCreated, "Successfully created a new plan", planObject)
	return c.JSON(response.StatusCode, response)
}

// Update plan
func UpdatePlanController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var updatedPlan models.Plan

	planID := c.Param("id")

	var existingPlan models.Plan
	existingPlan.InsertID(planID, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = c.Bind(&updatedPlan)
	if err.IsError() {
		err.StatusCode = http.StatusBadRequest
		err.ErrorReason = "Invalid request body"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetPlan(&existingPlan, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	if updatedPlan.Name != "" {
		existingPlan.Name = updatedPlan.Name
	}

	if updatedPlan.Duration != 0 {
		existingPlan.Duration = updatedPlan.Duration
	}

	if updatedPlan.Price != 0 {
		existingPlan.Price = updatedPlan.Price
	}

	if updatedPlan.Description != "" {
		existingPlan.Description = updatedPlan.Description
	}

	database.UpdatePlan(&existingPlan, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success(http.StatusOK, "Successfully updated plan", existingPlan)
	return c.JSON(response.StatusCode, response)
}

// Delete plan
func DeletePlanController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var planObject models.Plan

	planObject.InsertID(c.Param("id"), &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetPlan(&planObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.DeletePlan(&planObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	deletedPlan := map[string]int{
		"plan_id": int(planObject.ID),
	}
	response.Success(http.StatusCreated, "success delete plan", deletedPlan)
	return c.JSON(http.StatusOK, response)
}
