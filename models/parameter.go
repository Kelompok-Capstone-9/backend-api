package models

import (
	"gofit-api/constants"
	"strconv"
)

type Pages struct {
	PageString string
	Page       int
	Offset     int
	Limit      int
}

func (p *Pages) ConvertPageToINT(err *CustomError) {
	if p.PageString != "" {
		p.Page, err.ErrorMessage = strconv.Atoi(p.PageString)
		if err.ErrorMessage != nil {
			err.StatusCode = 400
			err.ErrorReason = "invalid page parameter :" + p.PageString
		}
	} else {
		p.Page = 1
	}
}

func (p *Pages) CalcOffsetLimit() (int, int) {
	var offset int
	if p.Page > 1 {
		offset = p.Page * constants.LIMIT
	}
	return offset, constants.LIMIT
}
