package controllers

import (
	"gofit-api/lib/database"
	"gofit-api/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetClassPackagesController(c echo.Context) error {
	var response models.GeneralListResponse
	var params models.ClassPackageParameters
	var page models.Pages
	var classPackages []models.ReadableClassPackage
	var totalData int
	var err models.CustomError

	// pagination
	page.PageString = c.QueryParam("page")
	page.PageSizeString = c.QueryParam("page_size")
	page.Paginate(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	//quey params
	params.PeriodString = c.QueryParam("period")
	params.MinPriceString = c.QueryParam("min_price")
	params.MaxPriceString = c.QueryParam("max_price")
	params.ClassIDString = c.QueryParam("class_id")

	params.ConvertAllParamStringToParams(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	if params.ParamIsSet() {
		query := params.DecodeToQueryString(&err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}

		classPackages, response.DataShown = database.GetClassPackagesWithParams(query, &page, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	} else {
		classPackages, response.DataShown = database.GetClassPackages(&page, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	}

	totalData = database.ClassPackageTotalData()

	response.Success("success get class packages", page.Page, totalData, classPackages)
	return c.JSON(response.StatusCode, response)
}

// get class by id
func GetClassPackageByIDController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var readableClassPackage models.ReadableClassPackage
	var classPackageObject models.ClassPackage

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classPackageObject.ID = uint(idParam.ID)

	database.GetClassPackage(&classPackageObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classPackageObject.ToReadableClassPackage(&readableClassPackage)

	response.Success(http.StatusOK, "success get class package", readableClassPackage)
	return c.JSON(response.StatusCode, response)
}

// create new class
func CreateClassPackageController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableClassPackage models.ReadableClassPackage
	var classPackageObject models.ClassPackage
	var classObject models.Class

	err.ErrorMessage = c.Bind(&readableClassPackage)
	if err.IsError() {
		err.ErrBind("invalid body request")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// validate class field
	err.ErrorMessage = readableClassPackage.Validate()
	if err.IsError() {
		err.ErrValidate("invalid field. field cant be blank")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classObject.ID = uint(readableClassPackage.Class.ID)
	database.GetClass(&classObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		response.ErrorReason = "invalid class"
		return c.JSON(response.StatusCode, response)
	}

	readableClassPackage.ToClassPackageObject(&classPackageObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classPackageObject.Class = classObject

	database.CreateClassPackage(&classPackageObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classPackageObject.ToReadableClassPackage(&readableClassPackage)

	response.Success(http.StatusCreated, "success create new class package", readableClassPackage)
	return c.JSON(response.StatusCode, response)
}

// edit class by id
func EditClassPackageController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var readableModifiedClassPackage models.ReadableClassPackage
	var readableClassPackage models.ReadableClassPackage
	var classPackageObject models.ClassPackage

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classPackageObject.ID = uint(idParam.ID)

	err.ErrorMessage = c.Bind(&readableModifiedClassPackage)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid body request"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = readableModifiedClassPackage.EditValidate()
	if err.IsError() {
		err.ErrValidate("field cant be blank, atleast one field need to be fill")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetClassPackage(&classPackageObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	classPackageObject.ToReadableClassPackage(&readableClassPackage)

	//replace exist data with new one
	if readableModifiedClassPackage.Class.ID != 0 {
		classObject := models.Class{ID: uint(readableModifiedClassPackage.Class.ID)}
		database.GetClass(&classObject, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			response.ErrorReason = "invalid class"
			return c.JSON(response.StatusCode, response)
		}
		classObject.ToReadableClassOnly(&readableClassPackage.Class)
	}
	if readableModifiedClassPackage.Period != "" {
		readableClassPackage.Period = readableModifiedClassPackage.Period
	}
	if readableModifiedClassPackage.Price != 0 {
		readableClassPackage.Price = readableModifiedClassPackage.Price
	}

	readableClassPackage.ToClassPackageObject(&classPackageObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.UpdateClassPackage(&classPackageObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success(http.StatusOK, "success edit class package", readableClassPackage)
	return c.JSON(http.StatusOK, response)
}

func DeleteClassPackageController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var classPackageObject models.ClassPackage

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	classPackageObject.ID = uint(idParam.ID)
	database.GetClassPackage(&classPackageObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.DeleteClassPackage(&classPackageObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	deletedClass := map[string]int{
		"class_package_id": int(classPackageObject.ID),
	}
	response.Success(http.StatusOK, "success delete class package", deletedClass)
	return c.JSON(http.StatusOK, response)
}
