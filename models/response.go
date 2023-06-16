package models

import (
	"gofit-api/constants"
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

	LoginResponse struct {
		ResponseMetadata `json:"metadata"`
		Data             interface{} `json:"data"`
		Token            string      `json:"token"`
	}
)

func (p *Pagination) CalculateShownData() {
	if p.TotalData < constants.LIMIT {
		p.DataShown = p.TotalData
	} else {
		p.DataShown = constants.LIMIT
	}
}

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
	glr.CalculateShownData()
}

func (gr *GeneralResponse) Success(statusCode int, message string, data interface{}) {
	gr.StatusCode = statusCode
	gr.Message = message
	gr.Data = data
}

func (lr *LoginResponse) Success(message string, data interface{}, token string) {
	lr.StatusCode = http.StatusOK
	lr.Message = message
	lr.Data = data
	lr.Token = token
}
