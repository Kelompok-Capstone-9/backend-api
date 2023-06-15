package controllers

import (
	"gofit-api/lib/database"
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