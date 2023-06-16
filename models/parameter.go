package models

import (
	"fmt"
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

type ClassParameters struct {
	ClassName string
	LocationParameters
	Date      string
	Time      string
	ClassType string
}

func (cp *ClassParameters) QueryString() string {
	var queryString string
	if cp.ClassName != "" {
		cp.ClassName = "%" + cp.ClassName + "%"
		switch queryString {
		case "":
			queryString += "name LIKE " + cp.ClassName
		default:
			queryString += "AND name LIKE " + cp.ClassName
		}
	}

	if cp.ClassName != "" {
		cp.ClassName = "%" + cp.ClassName + "%"
		switch queryString {
		case "":
			queryString += "name LIKE " + cp.ClassName
		default:
			queryString += "AND name LIKE " + cp.ClassName
		}
	}
	return queryString
}

type LocationParameters struct {
	Name    string
	Address string
	City    string
}

func (lp *LocationParameters) QueryString() string {
	var queryString string
	if lp.Name != "" {
		lp.Name = "%" + lp.Name + "%"
		switch queryString {
		case "":
			queryString += "name LIKE " + lp.Name
		default:
			queryString += "AND name LIKE " + lp.Name
		}
	}
	return queryString
}

type ClassPackageParameters struct {
	PeriodString   string
	Period         ClassPeriod
	MaxPriceString string
	MaxPrice       int
	MinPriceString string
	MinPrice       int
	ClassIDString  string
	ClassID        int
}

func (cpp *ClassPackageParameters) ParamIsSet() bool {
	var isSet bool

	if cpp.ClassID != 0 || cpp.Period != "" || cpp.MinPrice != 0 || cpp.MaxPrice != 0 {
		isSet = true
	}

	return isSet
}

func (cpp *ClassPackageParameters) ConvertAllParamStringToParams(err *CustomError) {
	cpp.ConvertPeriodStringToPeriod(err)
	cpp.ConvertPriceStringToInt(err)
	cpp.ConvertClassIDStringToINT(err)
}

func (cpp *ClassPackageParameters) DecodeToQueryString(err *CustomError) string {
	var queryString string
	isFilled := func() bool {
		return len(queryString) > 0
	}

	if cpp.Period != "" {
		if isFilled() {
			queryString += fmt.Sprintf(" AND period = '%s'", cpp.Period)
		} else {
			queryString += fmt.Sprintf("period = '%s'", cpp.Period)
		}
	}

	if cpp.MinPrice != 0 || cpp.MaxPrice != 0 {
		if isFilled() {
			queryString += fmt.Sprintf(" AND price BETWEEN %d AND %d", cpp.MinPrice, cpp.MaxPrice)
		} else {
			queryString += fmt.Sprintf("price BETWEEN %d AND %d", cpp.MinPrice, cpp.MaxPrice)
		}
	}

	if cpp.ClassID != 0 {
		if isFilled() {
			queryString += fmt.Sprintf(" AND class_id = %d", cpp.ClassID)
		} else {
			queryString += fmt.Sprintf("class_id = %d", cpp.ClassID)
		}
	}
	return queryString
}

func (cpp *ClassPackageParameters) ConvertPeriodStringToPeriod(err *CustomError) {
	if cpp.PeriodString != "" {
		cpp.Period, err.ErrorMessage = GenerateClassPeriod(cpp.PeriodString)
		if err.ErrorMessage != nil {
			err.StatusCode = 400
			err.ErrorReason = "invalid period parameter :" + cpp.PeriodString
		}
	}
}

func (cpp *ClassPackageParameters) ConvertPriceStringToInt(err *CustomError) {
	if cpp.MinPriceString != "" {
		cpp.MinPrice, err.ErrorMessage = strconv.Atoi(cpp.MinPriceString)
		if err.ErrorMessage != nil {
			err.StatusCode = 400
			err.ErrorReason = "invalid min price parameter :" + cpp.MinPriceString
		}
	}

	if cpp.MaxPriceString != "" {
		cpp.MaxPrice, err.ErrorMessage = strconv.Atoi(cpp.MaxPriceString)
		if err.ErrorMessage != nil {
			err.StatusCode = 400
			err.ErrorReason = "invalid max price parameter :" + cpp.MaxPriceString
		}
	}
}

func (cpp *ClassPackageParameters) ConvertClassIDStringToINT(err *CustomError) {
	if cpp.ClassIDString != "" {
		cpp.ClassID, err.ErrorMessage = strconv.Atoi(cpp.ClassIDString)
		if err.ErrorMessage != nil {
			err.StatusCode = 400
			err.ErrorReason = "invalid class id parameter :" + cpp.ClassIDString
		}
	}
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
