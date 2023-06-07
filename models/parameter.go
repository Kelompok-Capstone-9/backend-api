package models

import (
	"gofit-api/constants"
	"strconv"
)

type GeneralParameter struct {
	Name string
	Page Pages
}

func (gp *GeneralParameter) NameQueryForm() {
	gp.Name = "%%" + gp.Name + "%%"
}

type IDParameter struct {
	IDString string
	ID       int
}

func (i *IDParameter) ConvertIDStringToINT(err *CustomError) {
	if i.IDString != "" {
		i.ID, err.ErrorMessage = strconv.Atoi(i.IDString)
		if err.ErrorMessage != nil {
			err.StatusCode = 400
			err.ErrorReason = "invalid id parameter :" + i.IDString
		}
	}
}

type Pages struct {
	PageString string
	Page       int
	Offset     int
	Limit      int
}

func (p *Pages) ConvertPageStringToINT(err *CustomError) {
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

func (p *Pages) CalcOffsetLimit() {
	p.Limit = constants.LIMIT
	if p.Page > 1 {
		p.Offset = p.Page * p.Limit
	}
}
