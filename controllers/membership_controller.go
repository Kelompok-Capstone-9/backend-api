package controllers

import (
	"gofit-api/lib/database"
	"gofit-api/models"
	"net/http"
	"reflect"

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

	membershipObject.ToReadableMembership(&readableMembership)

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

	readableMembership.ToMembershipObject(&membershipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.CreateMembership(&membershipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetMembership(&membershipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	membershipObject.ToReadableMembership(&readableMembership)

	response.Success(http.StatusCreated, "Successfully created a new membership", readableMembership)
	return c.JSON(response.StatusCode, response)
}

// Update membership
func UpdateMembershipController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableModifiedMembership models.ReadableMembership
	var readableMembership models.ReadableMembership
	var membershipObject models.Membership

	membershipObject.InsertID(c.Param("id"), &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = c.Bind(&readableModifiedMembership)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid body request"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetMembership(&membershipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	membershipObject.ToReadableMembership(&readableMembership)

	//replace exist data with new one
	var membershipPointer *models.ReadableMembership = &readableMembership
	var modifiedMembershipPointer *models.ReadableMembership = &readableModifiedMembership
	membershipVal := reflect.ValueOf(membershipPointer).Elem()
	membershipType := membershipVal.Type()

	editVal := reflect.ValueOf(modifiedMembershipPointer).Elem()

	for i := 0; i < membershipVal.NumField(); i++ {
		//skip ID, CreatedAt, UpdatedAt field to be edited
		switch membershipType.Field(i).Name {
		case "ID":
			continue
		case "CreatedAt":
			continue
		case "UpdatedAt":
			continue
		}

		editField := editVal.Field(i)
		isSet := editField.IsValid() && !editField.IsZero()
		if isSet {
			membershipVal.Field(i).Set(editVal.Field(i))
		}
	}

	readableMembership.ToMembershipObject(&membershipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.UpdateMembership(&membershipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success(http.StatusCreated, "success edit membership", readableMembership)
	return c.JSON(http.StatusOK, response)
}

// Delete Membership
func DeleteMembershipController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var membershipObject models.Membership

	membershipObject.InsertID(c.Param("id"), &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetMembership(&membershipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.DeleteMembership(&membershipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	deletedMembership := map[string]int{
		"membership_id": int(membershipObject.ID),
	}
	response.Success(http.StatusCreated, "success delete membership", deletedMembership)
	return c.JSON(http.StatusOK, response)
}
