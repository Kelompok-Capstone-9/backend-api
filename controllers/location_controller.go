package controllers

import (
	"gofit-api/lib/database"
	"gofit-api/models"

	"github.com/labstack/echo/v4"
)

func GetLocationsController(c echo.Context) error {
	var response models.GeneralListResponse
	var params models.GeneralParameter
	var locations []models.ReadableLocation
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
		// locations, totalData = database.GetClassesWithParam(&params, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	default:
		locations, totalData = database.GetLocations(params.Page.Offset, params.Page.Limit, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	}

	response.Success("success get locations", params.Page.Page, totalData, locations)
	return c.JSON(response.StatusCode, response)
}