package controllers

import (
	"errors"
	"fmt"
	"gofit-api/lib/database"
	"gofit-api/lib/payment"
	"gofit-api/middlewares"
	"gofit-api/models"
	"net/http"
	"reflect"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

// GetTransactionsController retrieves all transactions
func GetTransactionsController(c echo.Context) error {
	var response models.GeneralListResponse
	var param models.GeneralParameter
	var err models.CustomError
	var transactions []models.ReadableTransaction
	var totalData int

	param.Page.PageString = c.QueryParam("page")
	param.Page.PageSizeString = c.QueryParam("page_size")
	param.Page.Paginate(&err)
	if err.IsError() {
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	param.Name = c.QueryParam("name")
	switch {
	case param.Name != "":
		param.NameQueryForm() // change name paramater to query form e.g: andy to %andy%
		transactions, response.DataShown = database.GetTransactionByUserName(param.Name, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	default:
		transactions, response.DataShown = database.GetTransactions(&param.Page, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	}

	totalData = database.CountTotalData("transactions")

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

// CreateTransactionController creates a new transaction
func PayController(c echo.Context) error {
	var response models.GeneralResponse
	var err models.CustomError

	var classTicketObject models.ClassTicket
	var membershipObject models.Membership
	var paymentMethodObject models.PaymentMethod

	var paymentRequest models.PaymentRequest
	var chargeRequest models.ChargeAPIRequest
	var paymentResult coreapi.ChargeResponse

	transaction := database.GetTransactionByCode(c.Param("transaction_code"), &err)
	if err.IsError() {
		err.ErrorReason = "this transaction doesn't exist"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}
	chargeRequest.TransactionDetail.OrderID = transaction.TransactionCode

	err.ErrorMessage = c.Bind(&paymentRequest)
	if err.IsError() {
		err.StatusCode = http.StatusBadRequest
		err.ErrorReason = "Invalid request body"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	err.ErrorMessage = paymentRequest.Validate()
	if err.IsError() {
		err.StatusCode = http.StatusBadRequest
		err.ErrorReason = "invalid request field"
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	switch transaction.Product {
	case models.ClassProduct:
		classTicketObject.ID = uint(transaction.ProductID)
		database.GetClassTicket(&classTicketObject, &err)
		if err.IsError() {
			err.ErrorReason = "class ticket doesn't exist"
			return c.JSON(response.StatusCode, response)
		}
		item := models.ItemDetail{
			ID:       strconv.Itoa(int(classTicketObject.ID)),
			Price:    int64(classTicketObject.ClassPackage.Price),
			Quantity: 1,
			Name:     fmt.Sprintf("%s - %s", classTicketObject.ClassPackage.Class.Name, classTicketObject.ClassPackage.Period),
		}
		chargeRequest.ItemDetails = append(chargeRequest.ItemDetails, item)

		chargeRequest.TransactionDetail.GrossAmount = int64(classTicketObject.ClassPackage.Price)
		chargeRequest.CustomerDetails.FirstName = classTicketObject.User.Name
		chargeRequest.CustomerDetails.Email = classTicketObject.User.Email
	case models.MembershipProduct:
		membershipObject.ID = uint(transaction.ProductID)
		database.GetMembership(&membershipObject, &err)
		if err.IsError() {
			err.ErrorReason = "membership doesn't exist"
			return c.JSON(response.StatusCode, response)
		}
		item := models.ItemDetail{
			ID:       strconv.Itoa(int(membershipObject.ID)),
			Price:    int64(membershipObject.Plan.Price),
			Quantity: 1,
			Name:     fmt.Sprintf("%s - %d Days", membershipObject.Plan.Name, membershipObject.Plan.Duration),
		}
		chargeRequest.ItemDetails = append(chargeRequest.ItemDetails, item)

		chargeRequest.TransactionDetail.GrossAmount = int64(membershipObject.Plan.Price)
		chargeRequest.CustomerDetails.FirstName = membershipObject.User.Name
		chargeRequest.CustomerDetails.Email = membershipObject.User.Email
	}

	switch paymentRequest.ReadablePaymentMethod.Name {
	case "gopay":
		paymentMethodObject.Name = "gopay"
		database.FirstOrCreatePaymentMethod(&paymentMethodObject, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
		var paymentMethodID *int = new(int)
		*paymentMethodID = int(paymentMethodObject.ID)
		transaction.PaymentMethodID = paymentMethodID
		transaction.PaymentMethod = paymentMethodObject
		chargeRequest.PaymentType = "gopay"
		err.StatusCode = http.StatusBadRequest
		err.ErrorMessage = errors.New("invalid payment method")
		err.ErrorReason = "sorry this payment method is not available yet."
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	case "credit_card":
		paymentMethodObject.Name = "gopay"
		database.FirstOrCreatePaymentMethod(&paymentMethodObject, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
		var paymentMethodID *int = new(int)
		*paymentMethodID = int(paymentMethodObject.ID)
		transaction.PaymentMethodID = paymentMethodID
		transaction.PaymentMethod = paymentMethodObject
		chargeRequest.PaymentType = "credit_card"
		chargeRequest.CreditCardToken, err.ErrorMessage = payment.GenerateCreditCardToken(&paymentRequest.CreditCard)
		if err.IsError() {
			err.StatusCode = http.StatusBadRequest
			err.ErrorReason = "Invalid Credit Card"
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
		var paymentError *midtrans.Error
		paymentResult, paymentError = payment.PayWithCreditCard(chargeRequest)
		if paymentError != nil {
			err.StatusCode = paymentError.StatusCode
			err.ErrorMessage = paymentError.RawError
			err.ErrorReason = paymentError.GetMessage()
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
	case "shopee_pay":
		paymentMethodObject.Name = "shoope_pay"
		database.FirstOrCreatePaymentMethod(&paymentMethodObject, &err)
		if err.IsError() {
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}
		var paymentMethodID *int = new(int)
		*paymentMethodID = int(paymentMethodObject.ID)
		transaction.PaymentMethodID = paymentMethodID
		transaction.PaymentMethod = paymentMethodObject
		chargeRequest.PaymentType = "shoope_pay"
		err.StatusCode = http.StatusBadRequest
		err.ErrorMessage = errors.New("invalid payment method")
		err.ErrorReason = "sorry this payment method is not available yet."
		response.ErrorOcurred(&err)
		return c.JSON(response.StatusCode, response)
	}

	switch paymentResult.TransactionStatus {
	case "capture", "settlement":
		// update transaction status to completed in database
		transaction.Status = "completed"
		database.UpdateTransaction(&transaction, &err)
		if err.IsError() {
			err.ErrorReason = "fail to update transaction status"
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}

		// update class / membership status in database
		switch transaction.Product {
		case models.ClassProduct:
			classTicketObject.Status = models.Booked
			database.UpdateClassTicket(&classTicketObject, &err)
			if err.IsError() {
				err.ErrorReason = "fail to update class ticket status"
				response.ErrorOcurred(&err)
				return c.JSON(response.StatusCode, response)
			}
		case models.MembershipProduct:
			membershipObject.IsActive = true
			database.UpdateMembership(&membershipObject, &err)
			if err.IsError() {
				err.ErrorReason = "fail to update membership status"
				response.ErrorOcurred(&err)
				return c.JSON(response.StatusCode, response)
			}
		}
	case "deny", "cancel", "expire":
		// update transaction status to cancel in database
		transaction.Status = "cancel"
		database.UpdateTransaction(&transaction, &err)
		if err.IsError() {
			err.ErrorReason = "fail to update transaction status"
			response.ErrorOcurred(&err)
			return c.JSON(response.StatusCode, response)
		}

		// update class / membership status in database
		switch transaction.Product {
		case models.ClassProduct:
			classTicketObject.Status = models.Cancelled
			database.UpdateClassTicket(&classTicketObject, &err)
			if err.IsError() {
				err.ErrorReason = "fail to update class ticket status"
				response.ErrorOcurred(&err)
				return c.JSON(response.StatusCode, response)
			}
		}

		response.StatusCode, _ = strconv.Atoi(paymentResult.StatusCode)
		response.Message = paymentResult.StatusMessage
		response.ErrorReason = "transaction failed try again later"
		return c.JSON(response.StatusCode, response)
	}

	var readableTransaction models.ReadableTransaction
	transaction.ToReadableTransaction(&readableTransaction)

	response.Success(http.StatusOK, paymentResult.StatusMessage, readableTransaction)
	return c.JSON(response.StatusCode, response)
}

func CallBackMidtrans(c echo.Context) error {
	var err models.CustomError
	var midtransCallback models.MidtransCallBack

	err.ErrorMessage = c.Bind(&midtransCallback)
	if err.IsError() {
		err.StatusCode = http.StatusBadRequest
		err.ErrorReason = "Invalid request body"
		return echo.NewHTTPError(err.StatusCode, fmt.Sprintf("%s:%s", err.ErrorMessage, err.ErrorReason))
	}

	switch midtransCallback.TransactionStatus {
	case "capture", "settlement":
		transaction := database.GetTransactionByCode(midtransCallback.TransactionCode, &err)
		if err.IsError() {
			fmt.Println(err)
			return echo.NewHTTPError(err.StatusCode, fmt.Sprintf("%s:%s", err.ErrorMessage, err.ErrorReason))
		}

		transaction.Status = "completed"
		database.UpdateTransaction(&transaction, &err)
		if err.IsError() {
			fmt.Println(err)
			return echo.NewHTTPError(err.StatusCode, fmt.Sprintf("%s:%s", err.ErrorMessage, err.ErrorReason))
		}

		switch transaction.Product {
		case models.MembershipProduct:
			err.ErrorMessage = database.ActivateMembershipByID(transaction.ProductID)
			if err.IsError() {
				fmt.Println(err)
				return echo.NewHTTPError(err.StatusCode, fmt.Sprintf("%s:%s", err.ErrorMessage, err.ErrorReason))
			}
			fmt.Println("membership(id=", transaction.ProductID, ") activated ")
		case models.ClassProduct:
			err.ErrorMessage = database.ChangeClassTicketStatus(transaction.ProductID, "booked")
			if err.IsError() {
				fmt.Println(err)
				return echo.NewHTTPError(err.StatusCode, fmt.Sprintf("%s:%s", err.ErrorMessage, err.ErrorReason))
			}
		}
	case "cancel", "expire":
		transaction := database.GetTransactionByCode(midtransCallback.TransactionCode, &err)
		if err.IsError() {
			fmt.Println(err)
			return echo.NewHTTPError(err.StatusCode, fmt.Sprintf("%s:%s", err.ErrorMessage, err.ErrorReason))
		}

		transaction.Status = "cancel"
		database.UpdateTransaction(&transaction, &err)
		if err.IsError() {
			fmt.Println(err)
			return echo.NewHTTPError(err.StatusCode, fmt.Sprintf("%s:%s", err.ErrorMessage, err.ErrorReason))
		}

		switch transaction.Product {
		case models.ClassProduct:
			err.ErrorMessage = database.ChangeClassTicketStatus(transaction.ProductID, "cancelled")
			if err.IsError() {
				fmt.Println(err)
				return echo.NewHTTPError(err.StatusCode, fmt.Sprintf("%s:%s", err.ErrorMessage, err.ErrorReason))
			}
		}
	}

	return c.HTML(200, "notification accepted")
}
