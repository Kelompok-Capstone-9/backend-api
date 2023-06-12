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

type ClassParameters struct{
	ClassName string
	LocationParameters
	Date string
	Time string
	ClassType string
}

func (cp *ClassParameters) QueryString() string {
	var queryString string
	if cp.ClassName != "" {
		cp.ClassName = "%" + cp.ClassName + "%"
		switch queryString{
		case "":
			queryString += "name LIKE " + cp.ClassName
		default:
			queryString += "AND name LIKE " + cp.ClassName
		}
	}

	if cp.ClassName != "" {
		cp.ClassName = "%" + cp.ClassName + "%"
		switch queryString{
		case "":
			queryString += "name LIKE " + cp.ClassName
		default:
			queryString += "AND name LIKE " + cp.ClassName
		}
	}
	return queryString
}

type LocationParameters struct{
	Name string
	Address string
	City string
}

func (lp *LocationParameters) QueryString() string {
	var queryString string
	if lp.Name != "" {
		lp.Name = "%" + lp.Name + "%"
		switch queryString{
		case "":
			queryString += "name LIKE " + lp.Name
		default:
			queryString += "AND name LIKE " + lp.Name
		}
	}
	return queryString
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
