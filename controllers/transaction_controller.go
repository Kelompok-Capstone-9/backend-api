package controllers

import (
	"gofit-api/lib/database"
	"gofit-api/middlewares"
	"gofit-api/models"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

// GetTransactionsController retrieves all transactions
func GetTransactionsController(c echo.Context) error {
	var response models.GeneralListResponse
	var param models.GeneralParameter
	var err models.CustomError
	var transactions []models.ReadableTransaction
	var totalData int

	param.Page.PageString = c.QueryParam("page")
	param.Page.ConvertPageStringToINT(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	param.Page.CalcOffsetLimit()
	param.Name = c.QueryParam("name")
	switch {
	case param.Name != "":
		param.NameQueryForm() // change name paramater to query form e.g: andy to %andy%
		transactions, totalData = database.GetTransactionByUserName(param.Name, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	default:
		transactions, totalData = database.GetTransactions(param.Page.Offset, param.Page.Limit, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	}

	response.Success("Successfully retrieved transactions", param.Page.Page, totalData, transactions)
	return c.JSON(http.StatusOK, response)
}

// // GetTransactionController retrieves a transaction by its ID
func GetTransactionController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableTransaction models.ReadableTransaction
	var transactionObject models.Transaction

	transactionID := c.Param("id")

	transactionObject.InsertID(transactionID, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetTransaction(&transactionObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	transactionObject.ToReadableTransaction(&readableTransaction)

	response.Success(http.StatusOK, "Successfully retrieved transaction", readableTransaction)
	return c.JSON(response.StatusCode, response)
}

// CreateTransactionController creates a new transaction
func CreateTransactionController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableTransaction models.ReadableTransaction
	var transactionObject models.Transaction

	err.ErrorMessage = c.Bind(&readableTransaction)
	if err.IsError() {
		err.StatusCode = http.StatusBadRequest
		err.ErrorReason = "Invalid request body"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	readableTransaction.ToTransactionObject(&transactionObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.CreateTransaction(&transactionObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetTransaction(&transactionObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	transactionObject.ToReadableTransaction(&readableTransaction)

	response.Success(http.StatusCreated, "Successfully created a new transaction", readableTransaction)
	return c.JSON(response.StatusCode, response)
}

// Get Transaction by user ID
func MyTransactionController(c echo.Context) error {
	UserID := middlewares.ExtractTokenUserID(c)
	var response models.GeneralResponse
	var err models.CustomError

	var readableTransaction models.ReadableTransaction
	var transactionObject models.Transaction

	database.GetTransactionByUserID(uint(UserID), &transactionObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	transactionObject.ToReadableTransaction(&readableTransaction)

	response.Success(http.StatusOK, "Successfully retrieved transaction", readableTransaction)
	return c.JSON(response.StatusCode, response)
}

// Update transaction
func UpdateTransactionController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var readableModifiedTransaction models.ReadableTransaction
	var readableTransaction models.ReadableTransaction
	var transactionObject models.Transaction

	transactionObject.InsertID(c.Param("id"), &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = c.Bind(&readableModifiedTransaction)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid body request"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetTransaction(&transactionObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	transactionObject.ToReadableTransaction(&readableTransaction)

	//replace exist data with new one
	var transactionPointer *models.ReadableTransaction = &readableTransaction
	var modifiedTransactionPointer *models.ReadableTransaction = &readableModifiedTransaction
	transactionVal := reflect.ValueOf(transactionPointer).Elem()
	transactionType := transactionVal.Type()

	editVal := reflect.ValueOf(modifiedTransactionPointer).Elem()

	for i := 0; i < transactionVal.NumField(); i++ {
		//skip ID, CreatedAt, UpdatedAt field to be edited
		switch transactionType.Field(i).Name {
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
			transactionVal.Field(i).Set(editVal.Field(i))
		}
	}

	readableTransaction.ToTransactionObject(&transactionObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.UpdateTransaction(&transactionObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	response.Success(http.StatusCreated, "success edit transaction", readableTransaction)
	return c.JSON(http.StatusOK, response)
}

// Delete Transaction
func DeleteTransactionController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var transactionObject models.Transaction

	transactionObject.InsertID(c.Param("id"), &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.GetTransaction(&transactionObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	database.DeleteTransaction(&transactionObject, &err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	deletedTransaction := map[string]int{
		"transaction_id": int(transactionObject.ID),
	}
	response.Success(http.StatusCreated, "success delete transaction", deletedTransaction)
	return c.JSON(http.StatusOK, response)
}
