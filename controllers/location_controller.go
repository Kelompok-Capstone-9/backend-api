package controllers

import (
	"gofit-api/lib/database"
	"gofit-api/models"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

func GetLocationsController(c echo.Context) error {
	var response models.GeneralListResponse
	var params models.LocationParameters
	var page models.Pages
	var locations []models.ReadableLocation
	var totalData int
	var err models.CustomError

	page.PageString = c.QueryParam("page")
	page.ConvertPageStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	page.CalcOffsetLimit()

	if params.ParamIsSet() {
		query := params.DecodeToQueryString()
		locations, totalData = database.GetLocationsWithParams(query, &page, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	} else {
		locations, totalData = database.GetLocations(page.Offset, page.Limit, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	}

	response.Success("success get locations", page.Page, totalData, locations)
	return c.JSON(response.StatusCode, response)
}

// get location by id
func GetLocationByIDController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var readableLocation models.ReadableLocation
	var locationObject models.Location

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	locationObject.ID = uint(idParam.ID)

	database.GetLocation(&locationObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	locationObject.ToReadableLocation(&readableLocation)

	response.Success(http.StatusOK, "success get location", readableLocation)
	return c.JSON(http.StatusOK, response)
}

// create new location
func CreateLocationController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableLocation models.ReadableLocation
	var locationObject models.Location

	err.ErrorMessage = c.Bind(&readableLocation)
	if err.IsError() {
		err.ErrBind("invalid body request")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// validate user field
	err.ErrorMessage = readableLocation.Validate()
	if err.IsError() {
		err.ErrValidate("invalid field. field cant be blank")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readableLocation.ToLocationObject(&locationObject)

	database.CreateLocation(&locationObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	locationObject.ToReadableLocation(&readableLocation)

	response.Success(http.StatusCreated, "success create new location", readableLocation)
	return c.JSON(response.StatusCode, response)
}

// update user by id
func EditLocationController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var readableModifiedLocation models.ReadableLocation
	var readableLocation models.ReadableLocation
	var locationObject models.Location

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	locationObject.ID = uint(idParam.ID)

	err.ErrorMessage = c.Bind(&readableModifiedLocation)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid body request"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetLocation(&locationObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	locationObject.ToReadableLocation(&readableLocation)

	//replace exist data with new one
	var locationPointer *models.ReadableLocation = &readableLocation
	var modifiedLocationPointer *models.ReadableLocation = &readableModifiedLocation
	locationVal := reflect.ValueOf(locationPointer).Elem()
	locationType := locationVal.Type()

	editVal := reflect.ValueOf(modifiedLocationPointer).Elem()

	for i := 0; i < locationVal.NumField(); i++ {
		//skip ID, CreatedAt, UpdatedAt field to be edited
		switch locationType.Field(i).Name {
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
			locationVal.Field(i).Set(editVal.Field(i))
		}
	}

	readableLocation.ToLocationObject(&locationObject)

	database.UpdateLocation(&locationObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success(http.StatusCreated, "success edit location", readableLocation)
	return c.JSON(http.StatusOK, response)
}

// delete location by id
func DeleteLocationController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var locationObject models.Location

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	locationObject.ID = uint(idParam.ID)

	database.GetLocation(&locationObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.DeleteLocation(&locationObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	deletedUser := map[string]int{
		"location_id": int(locationObject.ID),
	}
	response.Success(http.StatusCreated, "success delete location", deletedUser)
	return c.JSON(http.StatusOK, response)
}
