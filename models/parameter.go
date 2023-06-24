package models

import (
	"fmt"
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

func (lp *LocationParameters) ParamIsSet() bool {
	isSet := lp.Name != "" || lp.Address != "" || lp.City != ""
	return isSet
}

func (lp *LocationParameters) DecodeToQueryString() string {
	var queryString string
	isFilled := func() bool {
		return len(queryString) > 0
	}

	if lp.Name != "" {
		if isFilled() {
			queryString += fmt.Sprintf(" AND name = '%s'", lp.Name)
		} else {
			queryString += fmt.Sprintf("name = '%s'", lp.Name)
		}
	}
	if lp.Address != "" {
		if isFilled() {
			queryString += fmt.Sprintf(" AND address = '%s'", lp.Address)
		} else {
			queryString += fmt.Sprintf("address = '%s'", lp.Address)
		}
	}
	if lp.City != "" {
		if isFilled() {
			queryString += fmt.Sprintf(" AND city = '%s'", lp.City)
		} else {
			queryString += fmt.Sprintf("city = '%s'", lp.City)
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
	PageString     string
	Page           int
	PageSizeString string
	PageSize       int
	Offset         int
}

func (p *Pages) Paginate(err *CustomError) {
	if p.PageString != "" {
		p.Page, err.ErrorMessage = strconv.Atoi(p.PageString)
		if err.ErrorMessage != nil {
			err.StatusCode = 400
			err.ErrorReason = "invalid page parameter :" + p.PageString
		}
	}
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.PageSizeString != "" {
		p.PageSize, err.ErrorMessage = strconv.Atoi(p.PageSizeString)
		if err.ErrorMessage != nil {
			err.StatusCode = 400
			err.ErrorReason = "invalid page size parameter :" + p.PageSizeString
		}
	}

	switch {
	case p.PageSize > 100:
		p.PageSize = 100
	case p.PageSize <= 0:
		p.PageSize = 10
	}

	p.Offset = (p.Page - 1) * p.PageSize

}
