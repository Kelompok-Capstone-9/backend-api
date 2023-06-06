package controllers

import (
	"gofit-api/lib/database"
	"gofit-api/models"

	"github.com/labstack/echo/v4"
)

func GetClassesController(c echo.Context) error {
	var response models.GeneralListResponse
	var params models.GeneralParameter
	var classes []models.ReadableClass
	var totalData int
	var err models.CustomError

	params.Page.PageString = c.QueryParam("page")
	params.Page.ConvertPageStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	params.Page.CalcOffsetLimit()

	params.Name = c.QueryParam("name")
	switch {
	case params.Name != "":
		params.NameQueryForm() // change name paramater to query form e.g: andy to %andy%
		// classes, totalData = database.GetClassesWithParam(&params, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	default:
		classes, totalData = database.GetClasses(params.Page.Offset, params.Page.Limit, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	}

	response.Success("success get classes", params.Page.Page, totalData, classes)
	return c.JSON(response.StatusCode, response)
}

// func GetInstuctorController(c echo.Context) error {
// 	var response models.GeneralResponse
// 	var err models.CustomError
// 	var idParam models.IDParameter

// 	var ReadableInstructor models.ReadableInstructor
// 	var instructorObject models.Instructor

// 	idParam.IDString = c.Param("id")
// 	idParam.ConvertIDStringToINT(&err)
// 	if err.IsError() {
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}

// 	instructorObject.ID = uint(idParam.ID)
// 	database.GetInstructor(&instructorObject, &err)
// 	if err.IsError() {
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}

// 	instructorObject.ToReadableInstructor(&ReadableInstructor)

// 	response.Success(http.StatusOK, "success get instructor", ReadableInstructor)
// 	return c.JSON(response.StatusCode, response)
// }

// // create new instructor
// func CreateInstructorController(c echo.Context) error {
// 	var response models.GeneralResponse
// 	var err models.CustomError

// 	var readableInstructor models.ReadableInstructor
// 	var instructorObject models.Instructor

// 	err.ErrorMessage = c.Bind(&readableInstructor)
// 	if err.IsError() {
// 		err.ErrBind("invalid body request")
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}

// 	// validate instructor field
// 	err.ErrorMessage = readableInstructor.Validate()
// 	if err.IsError() {
// 		err.ErrValidate("invalid name or description. field cant be blank")
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}

// 	readableInstructor.ToInstructorObject(&instructorObject)

// 	database.CreateInstructor(&instructorObject, &err)
// 	if err.IsError() {
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}

// 	instructorObject.ToReadableInstructor(&readableInstructor)

// 	response.Success(http.StatusCreated, "success create new instructor", readableInstructor)
// 	return c.JSON(response.StatusCode, response)
// }

// func EditInstructorController(c echo.Context) error {
// 	var response models.GeneralResponse
// 	var err models.CustomError
// 	var idParam models.IDParameter

// 	var readableModifiedInstructor models.ReadableInstructor
// 	var readableInstructor models.ReadableInstructor
// 	var instructorObject models.Instructor

// 	idParam.IDString = c.Param("id")
// 	idParam.ConvertIDStringToINT(&err)
// 	if err.IsError() {
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}

// 	err.ErrorMessage = c.Bind(&readableModifiedInstructor)
// 	if err.IsError() {
// 		err.StatusCode = 400
// 		err.ErrorReason = "invalid body request"
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}

// 	err.ErrorMessage = readableModifiedInstructor.EditValidate()
// 	if err.IsError() {
// 		err.ErrValidate("field cant be blank, atleast one field need to be fill")
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}

// 	instructorObject.ID = uint(idParam.ID)
// 	database.GetInstructor(&instructorObject, &err)
// 	if err.IsError() {
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}
// 	instructorObject.ToReadableInstructor(&readableInstructor)

// 	if readableModifiedInstructor.Name != "" {
// 		readableInstructor.Name = readableModifiedInstructor.Name
// 	}
// 	if readableModifiedInstructor.Description != "" {
// 		readableInstructor.Description = readableModifiedInstructor.Description
// 	}

// 	readableInstructor.ToInstructorObject(&instructorObject)
// 	database.UpdateInstructor(&instructorObject, &err)
// 	if err.IsError() {
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}

// 	response.Success(http.StatusCreated, "success edit user", readableInstructor)
// 	return c.JSON(http.StatusOK, response)
// }

// func DeleteInstructorController(c echo.Context) error {
// 	var response models.GeneralResponse
// 	var err models.CustomError
// 	var idParam models.IDParameter

// 	var instructorObject models.Instructor

// 	idParam.IDString = c.Param("id")
// 	idParam.ConvertIDStringToINT(&err)
// 	if err.IsError() {
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}

// 	instructorObject.ID = uint(idParam.ID)
// 	database.GetInstructor(&instructorObject, &err)
// 	if err.IsError() {
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}

// 	database.DeleteInstructor(&instructorObject, &err)
// 	if err.IsError() {
// 		response.ErrorOcurred(&err)
// 		return c.JSON(response.StatusCode, response)
// 	}

// 	deletedUser := map[string]int{
// 		"user_id": int(instructorObject.ID),
// 	}
// 	response.Success(http.StatusCreated, "success delete user", deletedUser)
// 	return c.JSON(http.StatusOK, response)
// }
