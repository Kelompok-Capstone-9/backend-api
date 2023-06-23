package controllers

import (
	"gofit-api/lib/database"
	"gofit-api/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetPaymentMethodsController retrieves all payment method
func GetPaymentMethodsController(c echo.Context) error {
	var response models.GeneralListResponse
	var page models.Pages
	var err models.CustomError

	page.PageString = c.QueryParam("page")
	page.PageSizeString = c.QueryParam("page_size")
	page.Paginate(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	payment_methods, totalData := database.GetPaymentMethods(&page, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success("Successfully retrieved payment_method", page.Page, totalData, payment_methods)
	return c.JSON(http.StatusOK, response)
}

// Get payment method by ID
func GetPaymentMethodController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readablePaymentMethod models.ReadablePaymentMethod
	var paymentMethodObject models.PaymentMethod

	// Extract the payment ID from the URL parameter
	paymentMethodID := c.Param("id")

	// Convert the payment ID string to int
	paymentMethodObject.InsertID(paymentMethodID, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// Retrieve the payment method from the database
	database.GetPaymentMethod(&paymentMethodObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	// Convert the payment method object to a readable format
	paymentMethodObject.ToReadablePaymentMethod(&readablePaymentMethod)

	response.Success(http.StatusOK, "Successfully retrieved payment", readablePaymentMethod)
	return c.JSON(http.StatusOK, response)
}

// Create new payment
func CreatePaymentMethodController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readablePaymentMethod models.ReadablePaymentMethod
	var paymentMethodObject models.PaymentMethod

	err.ErrorMessage = c.Bind(&readablePaymentMethod)
	if err.IsError() {
		err.StatusCode = http.StatusBadRequest
		err.ErrorReason = "Invalid request body"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = readablePaymentMethod.Validate()
	if err.IsError() {
		err.ErrorReason = "invalid field"
		response.ErrorOcurred(&err)
	}

	readablePaymentMethod.ToPaymentMethodObject(&paymentMethodObject)

	database.CreatePaymentMethod(&paymentMethodObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success(http.StatusCreated, "Successfully created a new payment method", readablePaymentMethod)
	return c.JSON(response.StatusCode, response)
}

// Update payment method
func UpdatePaymentMethodController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var paymentMethodIDParam models.IDParameter

	var modifiedPaymentMethod models.ReadablePaymentMethod
	var readablePaymentMethod models.ReadablePaymentMethod
	var paymentMethodObject models.PaymentMethod

	paymentMethodIDParam.IDString = c.Param("id")
	paymentMethodIDParam.ConvertIDStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	paymentMethodObject.ID = uint(paymentMethodIDParam.ID)

	err.ErrorMessage = c.Bind(&modifiedPaymentMethod)
	if err.IsError() {
		err.StatusCode = http.StatusBadRequest
		err.ErrorReason = "Invalid request body"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = modifiedPaymentMethod.Validate()
	if err.IsError() {
		err.ErrorReason = "invalid field"
		response.ErrorOcurred(&err)
	}

	database.GetPaymentMethod(&paymentMethodObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	paymentMethodObject.Name = modifiedPaymentMethod.Name

	database.UpdatePaymentMethod(&paymentMethodObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	paymentMethodObject.ToReadablePaymentMethod(&readablePaymentMethod)

	response.Success(http.StatusOK, "Successfully updated payment method", readablePaymentMethod)
	return c.JSON(response.StatusCode, response)
}

// Delete payment method
func DeletePaymentMethodController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var paymentMethodObject models.PaymentMethod

	paymentMethodObject.InsertID(c.Param("id"), &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetPaymentMethod(&paymentMethodObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.DeletePaymentMethod(&paymentMethodObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	DeletePaymentMethod := map[string]int{
		"payment_method_id": int(paymentMethodObject.ID),
	}
	response.Success(http.StatusCreated, "success delete payment method", DeletePaymentMethod)
	return c.JSON(http.StatusOK, response)
}
