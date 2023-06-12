package controllers

import (
	"gofit-api/lib/database"
	"gofit-api/models"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

func GetClassPackgesController(c echo.Context) error {
	var response models.GeneralListResponse
	var params models.GeneralParameter
	var classPackages []models.ReadableClassPackage
	var totalData int
	var err models.CustomError

	params.Page.PageString = c.QueryParam("page")
	params.Page.ConvertPageStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	params.Page.CalcOffsetLimit()

	classPackages, totalData = database.GetClassPackages(params.Page.Offset, params.Page.Limit, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success("success get class packages", params.Page.Page, totalData, classPackages)
	return c.JSON(response.StatusCode, response)
}

// get class by id
func GetClassPackageByIDController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var readableClass models.ReadableClass
	var classObject models.Class

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classObject.ID = uint(idParam.ID)

	database.GetClass(&classObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classObject.ToReadableClass(&readableClass)

	response.Success(http.StatusOK, "success get instructor", readableClass)
	return c.JSON(response.StatusCode, response)
}

// create new class
func CreateClassPackageController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableClass models.ReadableClass
	var classObject models.Class
	var locationObject models.Location

	err.ErrorMessage = c.Bind(&readableClass)
	if err.IsError() {
		err.ErrBind("invalid body request")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// validate class field
	err.ErrorMessage = readableClass.Validate()
	if err.IsError() {
		err.ErrValidate("invalid field. field cant be blank")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	locationObject.ID = uint(readableClass.Location.ID)
	if locationObject.ID != 0 {
		database.GetLocation(&locationObject, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			response.ErrorReason = "invalid location"
			return c.JSON(response.StatusCode, response)
		}
	}

	readableClass.ToClassObject(&classObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classObject.Location = locationObject

	database.CreateClass(&classObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classObject.ToReadableClass(&readableClass)

	response.Success(http.StatusCreated, "success create new class", readableClass)
	return c.JSON(response.StatusCode, response)
}

// edit class by id
func EditClassPackageController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var readableModifiedClass models.ReadableClass
	var readableClass models.ReadableClass
	var classObject models.Class

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classObject.ID = uint(idParam.ID)

	err.ErrorMessage = c.Bind(&readableModifiedClass)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid body request"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = readableModifiedClass.EditValidate()
	if err.IsError() {
		err.ErrValidate("field cant be blank, atleast one field need to be fill")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetClass(&classObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classObject.ToReadableClass(&readableClass)

	//replace exist data with new one
	var classPointer *models.ReadableClass = &readableClass
	var modifiedClassPointer *models.ReadableClass = &readableModifiedClass
	classVal := reflect.ValueOf(classPointer).Elem()
	classType := classVal.Type()

	editVal := reflect.ValueOf(modifiedClassPointer).Elem()

	for i := 0; i < classVal.NumField(); i++ {
		//skip ID, CreatedAt, UpdatedAt field to be edited
		switch classType.Field(i).Name {
		case "ID":
			continue
		case "Location":
			if readableModifiedClass.Location.ID != 0 {
				locationObject := models.Location{ID: uint(readableModifiedClass.Location.ID)}
				database.GetLocation(&locationObject, &err)
				if err.IsError() {
					response.ErrorOcurred(&err)
					return c.JSON(response.StatusCode, response)
				}
				locationObject.ToReadableLocation(&readableModifiedClass.Location)
			}
		case "CreatedAt":
			continue
		case "UpdatedAt":
			continue
		}

		editField := editVal.Field(i)
		isSet := editField.IsValid() && !editField.IsZero()
		if isSet {
			classVal.Field(i).Set(editVal.Field(i))
		}
	}

	readableClass.ToClassObject(&classObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.UpdateClass(&classObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success(http.StatusOK, "success edit class", readableClass)
	return c.JSON(http.StatusOK, response)
}

func DeleteClassPackageController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var classObject models.Class

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classObject.ID = uint(idParam.ID)
	database.GetClass(&classObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.DeleteClass(&classObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	deletedClass := map[string]int{
		"class_id": int(classObject.ID),
	}
	response.Success(http.StatusOK, "success delete class", deletedClass)
	return c.JSON(http.StatusOK, response)
}
