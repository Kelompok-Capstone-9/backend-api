package controllers

import (
	"fmt"
	"gofit-api/lib/database"
	"gofit-api/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// get all healthtip tickets
func GetHealthtipsController(c echo.Context) error {
	var response models.GeneralListResponse
	var params models.GeneralParameter
	var healthtips []models.ReadableHealthtip
	var totalData int
	var err models.CustomError

	params.Page.PageString = c.QueryParam("page")
	params.Page.ConvertPageStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	params.Page.CalcOffsetLimit()

	healthtips, totalData = database.GetHealthtips(params.Page.Offset, params.Page.Limit, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success("success get healthtip", params.Page.Page, totalData, healthtips)
	return c.JSON(response.StatusCode, response)
}

// get healthtip by id
func GetHealthtipByIDController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var readableHealthtip models.ReadableHealthtip
	var healthtipObject models.Healthtip

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	healthtipObject.ID = uint(idParam.ID)

	database.GetHealthtip(&healthtipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	healthtipObject.ToReadableHealthtip(&readableHealthtip)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readableHealthtip.User.HidePassword()

	response.Success(http.StatusOK, "success get healthtip", readableHealthtip)
	return c.JSON(response.StatusCode, response)
}

// create new healthtip
func CreateHealthtipController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableHealthtip models.ReadableHealthtip
	var healthtipObject models.Healthtip
	var userObject models.User

	err.ErrorMessage = c.Bind(&readableHealthtip)
	if err.IsError() {
		err.ErrBind("invalid body request")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// validate healthtip field
	err.ErrorMessage = readableHealthtip.Validate()
	if err.IsError() {
		err.ErrValidate("invalid field. field cant be blank")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	userObject.ID = uint(readableHealthtip.User.ID)
	database.GetUser(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		response.ErrorReason = "invalid user"
		return c.JSON(response.StatusCode, response)
	}

	readableHealthtip.ToHealthtipObject(&healthtipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	healthtipObject.User = userObject

	database.CreateHealthtip(&healthtipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	healthtipObject.ToReadableHealthtip(&readableHealthtip)
	readableHealthtip.User.HidePassword()

	response.Success(http.StatusCreated, "success create new healthtip ticket", readableHealthtip)
	return c.JSON(response.StatusCode, response)
}

// edit healthtip ticket by id
func EditHealthtipController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var readableModifiedHealthtip models.ReadableHealthtip
	var readableHealthtip models.ReadableHealthtip
	var healthtipObject models.Healthtip

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	healthtipObject.ID = uint(idParam.ID)

	err.ErrorMessage = c.Bind(&readableModifiedHealthtip)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid body request"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = readableModifiedHealthtip.EditValidate()
	if err.IsError() {
		err.ErrValidate("field cant be blank, atleast one field need to be fill")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetHealthtip(&healthtipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	healthtipObject.ToReadableHealthtip(&readableHealthtip)

	//replace exist data with new one
	if readableModifiedHealthtip.User.ID != 0 {
		userObject := models.User{ID: uint(readableModifiedHealthtip.User.ID)}
		database.GetUser(&userObject, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			response.ErrorReason = "invalid user"
			return c.JSON(response.StatusCode, response)
		}
		userObject.ToReadableUser(&readableHealthtip.User)
	}
	if readableModifiedHealthtip.Title != "" {
		readableHealthtip.Title = readableModifiedHealthtip.Title
	}
	if readableModifiedHealthtip.Content != "" {
		readableHealthtip.Content = readableModifiedHealthtip.Content
	}
	if readableModifiedHealthtip.Image != "" {
		readableHealthtip.Image = readableModifiedHealthtip.Image
	}
	if readableModifiedHealthtip.PublishDate != "" {
		readableHealthtip.PublishDate = readableModifiedHealthtip.PublishDate
	}

	readableHealthtip.ToHealthtipObject(&healthtipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.UpdateHealthtip(&healthtipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success(http.StatusOK, "success edit healthtip", readableHealthtip)
	return c.JSON(http.StatusOK, response)
}

func DeleteHealthtipController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError
	var idParam models.IDParameter

	var healthtipObject models.Healthtip

	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	healthtipObject.ID = uint(idParam.ID)
	database.GetHealthtip(&healthtipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.DeleteHealthtip(&healthtipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	deletedHealthtip := map[string]int{
		"healthtip_id": int(healthtipObject.ID),
	}
	response.Success(http.StatusOK, "success delete healthtip", deletedHealthtip)
	return c.JSON(http.StatusOK, response)
}

func UploadHealthtipImageController(c echo.Context) error {
	var healthtipObject models.Healthtip
	var imageFile models.UploadImage
	var err models.CustomError

	var idParam models.IDParameter
	var response models.GeneralResponse

	// get user info
	idParam.IDString = c.Param("id")
	idParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	healthtipObject.ID = uint(idParam.ID)
	database.GetHealthtip(&healthtipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// image file
	imageFile.Name = fmt.Sprintf("healthtip%d", healthtipObject.ID)
	imageFile.Image, err.ErrorMessage = c.FormFile("file")
	if err.IsError() {
		err.NewError(http.StatusBadRequest, err.ErrorMessage, "invalid uploaded file")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = imageFile.Validate()
	if err.IsError() {
		err.NewError(http.StatusBadRequest, err.ErrorMessage, "invalid uploaded file")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = imageFile.VerifyImageExtension()
	if err.IsError() {
		err.NewError(http.StatusBadRequest, err.ErrorMessage, "invalid uploaded file")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	var imagePath string
	imagePath, err.ErrorMessage = imageFile.CopyIMGToAssets()
	if err.IsError() {
		err.NewError(http.StatusInternalServerError, err.ErrorMessage, "something went wrong when upload file")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	healthtipObject.Image = imagePath
	database.UpdateHealthtip(&healthtipObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	var readableHealthtip models.ReadableHealthtip
	healthtipObject.ToReadableHealthtip(&readableHealthtip)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success(http.StatusOK, "success upload healthtip image", readableHealthtip)
	return c.JSON(response.StatusCode, response)
}
