package controllers

import (
	"gofit-api/lib/database"
	"gofit-api/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetMembershipsController retrieves all memberships
func GetMembershipsController(c echo.Context) error {
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
	memberships, totalData := database.GetMemberships(offset, limit, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success("Successfully retrieved memberships", page.Page, totalData, memberships)
	return c.JSON(http.StatusOK, response)
}

// // GetMembershipController retrieves a membership by its ID
func GetMembershipController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableMembership models.ReadableMembership
	var membershipObject models.Membership

	membershipID := c.Param("id")

	membershipObject.InsertID(membershipID, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetMembership(&membershipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readableMembership.ToMembershipObject(&membershipObject)

	response.Success(http.StatusOK, "Successfully retrieved membership", readableMembership)
	return c.JSON(response.StatusCode, response)
}

// CreateMembershipController creates a new membership
func CreateMembershipController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableMembership models.ReadableMembership
	var membershipObject models.Membership

	err.ErrorMessage = c.Bind(&readableMembership)
	if err.IsError() {
		err.StatusCode = http.StatusBadRequest
		err.ErrorReason = "Invalid request body"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readableMembership.ToMembershipObject(&membershipObject)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.CreateMembership(&membershipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readableMembership.ToMembershipObject(&membershipObject)

	response.Success(http.StatusCreated, "Successfully created a new membership", readableMembership)
	return c.JSON(response.StatusCode, response)
}
