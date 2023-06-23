package models

import (
	"net/http"
)

type (
	ResponseMetadata struct {
		StatusCode  int    `json:"status_code"`
		Message     string `json:"message"`
		ErrorReason string `json:"error_reason"`
	}

	Pagination struct {
		Page      int `json:"page"`
		DataShown int `json:"data_shown"`
		TotalData int `json:"total_data"`
	}

	GeneralListResponse struct {
		ResponseMetadata `json:"metadata"`
		Pagination       `json:"pagination"`
		Data             interface{} `json:"data"`
	}

	GeneralResponse struct {
		ResponseMetadata `json:"metadata"`
		Data             interface{} `json:"data"`
	}

	ProductTransactionInfoResponse struct {
		TransactionCode string `json:"transaction_code"`
		Message         string `json:"message"`
		TransactionLink string `json:"transaction_link"`
	}

	ProductResponse struct {
		ResponseMetadata               `json:"metadata"`
		ProductTransactionInfoResponse `json:"transaction_info"`
		Data                           interface{} `json:"data"`
	}

	LoginResponse struct {
		ResponseMetadata `json:"metadata"`
		Data             interface{} `json:"data"`
		Token            string      `json:"token"`
	}
)

func (rm *ResponseMetadata) ErrorOcurred(err *CustomError) {
	rm.StatusCode = err.StatusCode
	rm.Message = err.ErrorMessage.Error()
	rm.ErrorReason = err.ErrorReason
}

func (glr *GeneralListResponse) Success(message string, page, totalData int, data interface{}) {
	glr.StatusCode = http.StatusOK
	glr.Message = message
	glr.Page = page
	glr.TotalData = totalData
	glr.Data = data
}

func (gr *GeneralResponse) Success(statusCode int, message string, data interface{}) {
	gr.StatusCode = statusCode
	gr.Message = message
	gr.Data = data
}

func (ptir *ProductTransactionInfoResponse) TransactionCreated(transactionCode string, transactionMessage string, transactionLink string) {
	ptir.TransactionCode = transactionCode
	ptir.Message = transactionMessage
	ptir.TransactionLink = transactionLink
}

func (ttr *ProductResponse) Success(statusCode int, message string, data interface{}) {
	ttr.StatusCode = statusCode
	ttr.ResponseMetadata.Message = message
	ttr.Data = data
}

func (lr *LoginResponse) Success(message string, data interface{}, token string) {
	lr.StatusCode = http.StatusOK
	lr.Message = message
	lr.Data = data
	lr.Token = token
}
