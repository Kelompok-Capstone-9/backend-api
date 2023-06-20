package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"reflect"

	"gofit-api/lib/database"
	"gofit-api/lib/mailer"
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
	page.ConvertPageStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	page.CalcOffsetLimit()
	users, totalData := database.GetUsers(page.Offset, page.Limit, &err)
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
	var passwordIsModified bool

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

	err.ErrorMessage = readableModifiedUser.EditValidate()
	if err.IsError() {
		response.ErrorOcurred(&err)
		response.ErrorReason = "invalid field"
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
		case "Password":
			if readableModifiedUser.Password != "" {
				passwordIsModified = true
			} else {
				continue
			}
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

	if passwordIsModified {
		userObject.HashingPassword(&err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	}

	database.UpdateUser(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readableUser.HidePassword()
	response.Success(http.StatusOK, "success edit user", readableUser)
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
	var membershipObject models.Membership

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

	database.GetMembershipByUserID(userObject.ID, &membershipObject, &err)
	if err.IsError() {
		if err.StatusCode == 500 {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	}

	var token string
	token, err.ErrorMessage = middlewares.CreateToken(int(userObject.ID), userObject.Email, membershipObject.CheckMembershipActivity(), userObject.IsAdmin)
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

func UploadProfilePictureController(c echo.Context) error {
	var userObject models.User
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

	userObject.ID = uint(idParam.ID)
	database.GetUser(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// image file
	imageFile.Name = fmt.Sprintf("user%d", userObject.ID)
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

	userObject.ProfilePicture = imagePath
	database.UpdateUser(&userObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	var readableUser models.ReadableUser
	userObject.ToReadableUser(&readableUser)
	readableUser.HidePassword()

	response.Success(http.StatusOK, "success upload user profile image", readableUser)
	return c.JSON(response.StatusCode, response)
}

func ForgotPasswordController(c echo.Context) error {
	var response models.ResponseMetadata
	var err models.CustomError
	var otp models.OTP

	err.ErrorMessage = c.Bind(&otp)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid body request"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	if otp.Email == "" {
		err.NewError(http.StatusBadRequest, errors.New("invalid email"), "email cant be blank")
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	} else {
		_, emailError := mail.ParseAddress(otp.Email)
		if emailError != nil {
			err.NewError(http.StatusBadRequest, errors.New("" + otp.Email + " is not an email"), "invalid email")
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	}

	var message string
	err.ErrorMessage = otp.GenerateOTP()
	if err.IsError() {
		err.StatusCode = http.StatusInternalServerError
		err.ErrorReason = "something went wrong when generation OTP"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	message, err.ErrorMessage = mailer.SendOTP(otp.Email, otp.Code)
	if err.IsError() {
		err.StatusCode = http.StatusInternalServerError
		err.ErrorReason = "something went wrong when sending email"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.StatusCode = http.StatusOK
	response.Message = message
	return c.JSON(response.StatusCode, response)
}
