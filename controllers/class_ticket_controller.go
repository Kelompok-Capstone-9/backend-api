package controllers

import (
	"errors"
	"fmt"
	"gofit-api/lib/database"
	"gofit-api/middlewares"
	"gofit-api/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// get all class tickets
func GetClassTicketsController(c echo.Context) error {
	var response models.GeneralListResponse
	var params models.GeneralParameter
	var classTickets []models.ReadableClassTicket
	var totalData int
	var err models.CustomError

	params.Page.PageString = c.QueryParam("page")
	params.Page.ConvertPageStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	params.Page.CalcOffsetLimit()

	classTickets, totalData = database.GetClassTickets(params.Page.Offset, params.Page.Limit, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success("success get class ticket", params.Page.Page, totalData, classTickets)
	return c.JSON(response.StatusCode, response)
}

// get class ticket by id
func GetClassTicketByIDController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var readableClassTicket models.ReadableClassTicket
	var classTicketObject models.ClassTicket

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classTicketObject.ID = uint(idParam.ID)

	database.GetClassTicket(&classTicketObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classTicketObject.ToReadableClassTicket(&readableClassTicket)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readableClassTicket.User.HidePassword()

	response.Success(http.StatusOK, "success get class ticket", readableClassTicket)
	return c.JSON(response.StatusCode, response)
}

// create new class ticket
func CreateClassTicketController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableClassTicket models.ReadableClassTicket
	var classTicketObject models.ClassTicket
	var userObject models.User
	var classPackageObject models.ClassPackage

	err.ErrorMessage = c.Bind(&readableClassTicket)
	if err.IsError() {
		err.ErrBind("invalid body request")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// validate class field
	err.ErrorMessage = readableClassTicket.Validate()
	if err.IsError() {
		err.ErrValidate("invalid field. field cant be blank")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	userObject.ID = uint(readableClassTicket.User.ID)
	database.GetUser(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		response.ErrorReason = "invalid user"
		return c.JSON(response.StatusCode, response)
	}

	classPackageObject.ID = uint(readableClassTicket.ClassPackage.ID)
	database.GetClassPackage(&classPackageObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		response.ErrorReason = "invalid class package"
		return c.JSON(response.StatusCode, response)
	}

	readableClassTicket.ToClassTicketObject(&classTicketObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classTicketObject.User = userObject
	classTicketObject.ClassPackage = classPackageObject

	database.CreateClassTicket(&classTicketObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classTicketObject.ToReadableClassTicket(&readableClassTicket)
	readableClassTicket.User.HidePassword()

	response.Success(http.StatusCreated, "success create new class ticket", readableClassTicket)
	return c.JSON(response.StatusCode, response)
}

// edit class ticket by id
func EditClassTicketController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var readableModifiedClassTicket models.ReadableClassTicket
	var readableClassTicket models.ReadableClassTicket
	var classTicketObject models.ClassTicket

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classTicketObject.ID = uint(idParam.ID)

	err.ErrorMessage = c.Bind(&readableModifiedClassTicket)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid body request"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = readableModifiedClassTicket.EditValidate()
	if err.IsError() {
		err.ErrValidate("field cant be blank, atleast one field need to be fill")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetClassTicket(&classTicketObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classTicketObject.ToReadableClassTicket(&readableClassTicket)

	//replace exist data with new one
	if readableModifiedClassTicket.User.ID != 0 {
		userObject := models.User{ID: uint(readableModifiedClassTicket.User.ID)}
		database.GetUser(&userObject, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			response.ErrorReason = "invalid user"
			return c.JSON(response.StatusCode, response)
		}
		userObject.ToReadableUser(&readableClassTicket.User)
	}
	if readableModifiedClassTicket.ClassPackage.ID != 0 {
		classPackageObject := models.ClassPackage{ID: uint(readableModifiedClassTicket.ClassPackage.ID)}
		database.GetClassPackage(&classPackageObject, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			response.ErrorReason = "invalid class package"
			return c.JSON(response.StatusCode, response)
		}
		classPackageObject.ToReadableClassPackage(&readableClassTicket.ClassPackage)
	}
	if readableModifiedClassTicket.Status != "" {
		readableClassTicket.Status = readableModifiedClassTicket.Status
	}

	readableClassTicket.ToClassTicketObject(&classTicketObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.UpdateClassTicket(&classTicketObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success(http.StatusOK, "success edit class ticket", readableClassTicket)
	return c.JSON(http.StatusOK, response)
}

func DeleteClassTicketController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var classTicketObject models.ClassTicket

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classTicketObject.ID = uint(idParam.ID)
	database.GetClassTicket(&classTicketObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.DeleteClassTicket(&classTicketObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	deletedClass := map[string]int{
		"class_ticket_id": int(classTicketObject.ID),
	}
	response.Success(http.StatusOK, "success delete class ticket", deletedClass)
	return c.JSON(http.StatusOK, response)
}

// create new class ticket for users
func CreateMyTicketController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableClassTicket models.ReadableClassTicket
	var classTicketObject models.ClassTicket
	// var userObject models.User
	// var classPackageObject models.ClassPackage

	var classPackageIDParam models.IDParameter
	classPackageIDParam.IDString = c.Param("class_package_id")
	classPackageIDParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		response.ErrorReason = "invalid class package id"
		return c.JSON(response.StatusCode, response)
	}
	classTicketObject.ClassPackage.ID = uint(classPackageIDParam.ID)
	database.GetClassPackage(&classTicketObject.ClassPackage, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		response.ErrorReason = "invalid class package"
		return c.JSON(response.StatusCode, response)
	}

	userID := uint(middlewares.ExtractTokenUserID(c))
	classTicketObject.User.ID = userID
	database.GetUser(&classTicketObject.User, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		response.ErrorReason = "invalid user"
		return c.JSON(response.StatusCode, response)
	}

	database.CreateClassTicket(&classTicketObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classTicketObject.ToReadableClassTicket(&readableClassTicket)
	readableClassTicket.User.HidePassword()
	readableClassTicket.ClassPackage.Class.HideLink()

	response.Success(http.StatusCreated, "success create my class ticket. proceed to transaction to book", readableClassTicket)
	return c.JSON(response.StatusCode, response)
}

func GetMyTicketsController(c echo.Context) error {
	var response models.GeneralListResponse
	var page models.Pages
	var classTickets []models.ReadableClassTicket
	var totalData int
	var err models.CustomError

	page.PageString = c.QueryParam("page")
	page.ConvertPageStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	page.CalcOffsetLimit()

	userID := int(middlewares.ExtractTokenUserID(c))
	query := fmt.Sprintf("user_id = %d", userID)

	classTickets, totalData = database.GetClassTicketsWithParams(query, &page, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success("success get my class tickets", page.Page, totalData, classTickets)
	return c.JSON(response.StatusCode, response)
}

func GetMyTicketDetailController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var readableClassTicket models.ReadableClassTicket
	var classTicketObject models.ClassTicket

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classTicketObject.ID = uint(idParam.ID)

	database.GetClassTicket(&classTicketObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classTicketObject.ToReadableClassTicket(&readableClassTicket)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// check if the ticket belongs to the same user retrieve this data
	userID := middlewares.ExtractTokenUserID(c)
	isMyTicket := readableClassTicket.User.ID == int(userID)
	if !isMyTicket {
		err.NewError(http.StatusUnauthorized, errors.New("unauthorized"), "this ticket doesnt belong to you")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readableClassTicket.User.HidePassword()

	response.Success(http.StatusOK, "success get my class ticket detail", readableClassTicket)
	return c.JSON(response.StatusCode, response)
}

func CancelMyTicketController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var readableClassTicket models.ReadableClassTicket
	var classTicketObject models.ClassTicket

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classTicketObject.ID = uint(idParam.ID)

	database.GetClassTicket(&classTicketObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classTicketObject.Status = models.Cancelled
	database.UpdateClassTicket(&classTicketObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classTicketObject.ToReadableClassTicket(&readableClassTicket)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// check if the ticket belongs to the same user retrieve this data
	userID := middlewares.ExtractTokenUserID(c)
	isMyTicket := readableClassTicket.User.ID == int(userID)
	if !isMyTicket {
		err.NewError(http.StatusUnauthorized, errors.New("unauthorized"), "this ticket doesnt belong to you")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readableClassTicket.User.HidePassword()

	response.Success(http.StatusOK, "success cancel my class ticket", readableClassTicket)
	return c.JSON(response.StatusCode, response)
}
