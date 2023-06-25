package controllers

import (
	"fmt"
	"gofit-api/lib/database"
	"gofit-api/middlewares"
	"gofit-api/models"
	"net/http"
	"reflect"
	"time"

	"github.com/labstack/echo/v4"
)

// GetMembershipsController retrieves all memberships
func GetMembershipsController(c echo.Context) error {
	var response models.GeneralListResponse
	var param models.GeneralParameter
	var err models.CustomError
	var memberships []models.ReadableMembership
	var totalData int

	param.Page.PageString = c.QueryParam("page")
	param.Page.PageSizeString = c.QueryParam("page_size")
	param.Page.Paginate(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	param.Name = c.QueryParam("name")
	switch {
	case param.Name != "":
		param.NameQueryForm() // change name paramater to query form e.g: andy to %andy%
		memberships, response.DataShown = database.GetMembershipByUserName(param.Name, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	default:
		memberships, response.DataShown = database.GetMemberships(&param.Page, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	}

	totalData = database.CountTotalData("memberships")

	response.Success("Successfully retrieved memberships", param.Page.Page, totalData, memberships)
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

// Get Membership by user ID
func MyMembershipController(c echo.Context) error {
	UserID := middlewares.ExtractTokenUserID(c)
	var response models.GeneralResponse
	var err models.CustomError

	var readableMembership models.ReadableMembership
	var membershipObject models.Membership

	database.GetMembershipByUserID(uint(UserID), &membershipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	membershipObject.ToReadableMembership(&readableMembership)

	response.Success(http.StatusOK, "Successfully retrieved membership", readableMembership)
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

func JoinMembershipController(c echo.Context) error {
	var response models.ProductResponse
	var err models.CustomError

	var readableMembership models.ReadableMembership
	var membershipObject models.Membership
	var plan models.Plan

	userID := middlewares.ExtractTokenUserID(c)
	planID := models.IDParameter{IDString: c.Param("plan_id")}
	planID.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	plan.ID = uint(planID.ID)
	database.GetPlan(&plan, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		response.ErrorReason = "plan doesn't exists. make sure insert the right ID"
		return c.JSON(response.StatusCode, response)
	}

	membershipObject.UserID = uint(userID)
	membershipObject.PlanID = plan.ID
	membershipObject.StartDate = time.Now()
	membershipObject.EndDate = time.Now().AddDate(0, 0, plan.Duration)

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

	transactionObject := models.Transaction{
		TransactionCode: fmt.Sprintf("TM%d", membershipObject.ID),
		Product:         models.MembershipProduct,
		ProductID:       int(membershipObject.ID),
		Amount:          int(membershipObject.Plan.Price),
		Status:          string(models.Pending),
	}
	database.CreateTransaction(&transactionObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	transactionLink := fmt.Sprintf("/transactions/pay/%s", transactionObject.TransactionCode)
	response.TransactionCreated(transactionObject.TransactionCode, "transaction created continue to payment to activate your membership", transactionLink)

	membershipObject.ToReadableMembership(&readableMembership)

	response.Success(http.StatusCreated, "Successfully created a new membership", readableMembership)
	return c.JSON(response.StatusCode, response)
}
