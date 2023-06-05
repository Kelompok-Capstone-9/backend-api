package controllers

import (
	"net/http"
	"net/mail"
	"reflect"

	"gofit-api/lib/database"
	"gofit-api/middlewares"
	"gofit-api/models"

	"github.com/labstack/echo/v4"
)

// get all users
func GetUsersController(c echo.Context) error {
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
	users, totalData := database.GetUsers(offset, limit, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success("success get users", page.Page, totalData, users)
	return c.JSON(response.StatusCode, response)
}

// get user by id
func GetUserController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableUser models.ReadableUser
	var userObject models.User

	userObject.InsertID(c.Param("id"), &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetUser(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	userObject.ToReadableUser(&readableUser)
	readableUser.HidePassword()

	response.Success(http.StatusOK, "success get user", readableUser)
	return c.JSON(http.StatusOK, response)
}

// create new user
func CreateUserController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableUser models.ReadableUser
	var userObject models.User

	err.ErrorMessage = c.Bind(&readableUser)
	if err.IsError() {
		err.ErrBind("invalid body request")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// validate user field
	err.ErrorMessage = readableUser.Validate()
	if err.IsError() {
		err.ErrValidate("invalid field or email. field cant be blank or email must containt @email.com")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readableUser.ToUserObject(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.CreateUser(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	userObject.ToReadableUser(&readableUser)
	readableUser.HidePassword()

	response.Success(http.StatusCreated, "success create new user", readableUser)
	return c.JSON(response.StatusCode, response)
}

// update user by id
func EditUserController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableModifiedUser models.ReadableUser
	var readableUser models.ReadableUser
	var userObject models.User

	userObject.InsertID(c.Param("id"), &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = c.Bind(&readableModifiedUser)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid body request"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetUser(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	userObject.ToReadableUser(&readableUser)

	//replace exist data with new one
	var userPointer *models.ReadableUser = &readableUser
	var modifiedUserPointer *models.ReadableUser = &readableModifiedUser
	userVal := reflect.ValueOf(userPointer).Elem()
	userType := userVal.Type()

	editVal := reflect.ValueOf(modifiedUserPointer).Elem()

	for i := 0; i < userVal.NumField(); i++ {
		//skip ID, CreatedAt, UpdatedAt field to be edited
		switch userType.Field(i).Name {
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
			userVal.Field(i).Set(editVal.Field(i))
		}
	}

	readableUser.ToUserObject(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.UpdateUser(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readableUser.HidePassword()
	response.Success(http.StatusCreated, "success edit user", readableUser)
	return c.JSON(http.StatusOK, response)
}

// delete user by id
func DeleteUserController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var userObject models.User

	userObject.InsertID(c.Param("id"), &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetUser(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.DeleteUser(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	deletedUser := map[string]int{
		"user_id": int(userObject.ID),
	}
	response.Success(http.StatusCreated, "success delete user", deletedUser)
	return c.JSON(http.StatusOK, response)
}

func LoginUserController(c echo.Context) error {
	var response models.LoginResponse
	var err models.CustomError
	var loginReq models.LoginRequest

	err.ErrorMessage = c.Bind(&loginReq)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid body request"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// validate email
	_, emailError := mail.ParseAddress(loginReq.Email)
	if emailError != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "invalid email"
		response.ErrorReason = loginReq.Email + " is not an email"
		return c.JSON(response.StatusCode, response)
	}

	userObject := database.Login(loginReq.Email, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	match := userObject.MatchingPassword(loginReq.Password)
	if !match {
		err.FailLogin()
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	var token string
	token, err.ErrorMessage = middlewares.CreateToken(int(userObject.ID), userObject.Email, userObject.IsAdmin)
	if err.IsError() {
		err.StatusCode = 500
		err.ErrorReason = "fail to create jwt token"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	var readableUser models.ReadableUser
	userObject.ToReadableUser(&readableUser)
	readableUser.HidePassword()

	response.Success("success login", readableUser, token)
	return c.JSON(response.StatusCode, response)
}
