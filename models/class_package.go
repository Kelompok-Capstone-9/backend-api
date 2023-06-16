package models

import (
	"errors"
)

type ClassPackage struct {
	ID           uint
	Period       ClassPeriod `gorm:"type:enum('daily','weekly','monthly')"`
	Price        float32
	ClassID      uint
	Class        Class
	ClassTickets []ClassTicket
	Metadata     `gorm:"embedded"`
}

func (cp *ClassPackage) ToReadableClassPackageOnly(readableClassPackageOnly *ReadableClassPackageOnly) {
	readableMetadata := cp.Metadata.ToReadableMetadata()

	readableClassPackageOnly.ID = int(cp.ID)
	readableClassPackageOnly.Period = string(cp.Period)
	readableClassPackageOnly.Price = cp.Price
	readableClassPackageOnly.ClassID = int(cp.ClassID)
	readableClassPackageOnly.ReadableMetadata = *readableMetadata

}

func (cp *ClassPackage) ToReadableClassPackage(readableClassPackage *ReadableClassPackage) {
	readableMetadata := cp.Metadata.ToReadableMetadata()

	cp.Class.ToReadableClassOnly(&readableClassPackage.Class)

	readableClassPackage.ID = int(cp.ID)
	readableClassPackage.Period = string(cp.Period)
	readableClassPackage.Price = cp.Price

	readableClassPackage.ReadableMetadata = *readableMetadata
}

type ReadableClassPackageOnly struct {
	ID               int     `json:"id"`
	Period           string  `json:"period"`
	Price            float32 `json:"price"`
	ClassID          int     `json:"class_id"`
	ReadableMetadata `json:"metadata"`
}

type ReadableClassPackage struct {
	ID               int               `json:"id"`
	Period           string            `json:"period"`
	Price            float32           `json:"price"`
	Class            ReadableClassOnly `json:"class"`
	ReadableMetadata `json:"metadata"`
}

func (rcpwc *ReadableClassPackage) Validate() error {
	switch {
	case rcpwc.Period == "":
		return errors.New("invalid period")
	case rcpwc.Price == 0:
		return errors.New("invalid price")
	case rcpwc.Class.ID == 0:
		return errors.New("invalid class id")
	}

	if rcpwc.Period != "daily" && rcpwc.Period != "weekly" && rcpwc.Period != "monthly" {
		return errors.New("invalid period. period must contain daily or weekly or monthly")
	}

	return nil
}

func (rcpwc *ReadableClassPackage) EditValidate() error {
	allFieldBlank := rcpwc.Period == "" && rcpwc.Price == 0 && rcpwc.Class.ID == 0
	if allFieldBlank {
		return errors.New("all field is blank. nothing to change")
	}

	if rcpwc.Period != "" {
		if rcpwc.Period != "daily" && rcpwc.Period != "weekly" && rcpwc.Period != "monthly" {
			return errors.New("invalid period. must containt daily or weekly or monthly")
		}
	}

	return nil
}

func (rcpwc *ReadableClassPackage) ToClassPackageObject(classPackageObject *ClassPackage, err *CustomError) {
	var classPeriod ClassPeriod
	classPeriod, err.ErrorMessage = GenerateClassPeriod(rcpwc.Period)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid class period input"
	}

	classPackageObject.ID = uint(rcpwc.ID)
	classPackageObject.Period = classPeriod
	classPackageObject.Price = rcpwc.Price
	classPackageObject.ClassID = uint(rcpwc.Class.ID)
	classPackageObject.Class.ID = uint(rcpwc.Class.ID)
}

func ToReadableClassPackageOnlyList(classPackageObjectList []ClassPackage, err *CustomError) []ReadableClassPackageOnly {
	readableClassPacakgeList := make([]ReadableClassPackageOnly, len(classPackageObjectList))

	for i, item := range classPackageObjectList {
		var readableClassPackage ReadableClassPackageOnly
		item.ToReadableClassPackageOnly(&readableClassPackage)
		if err.IsError() {
			err.StatusCode = 500
			err.ErrorReason = "fail to parse class package"
			return nil
		}
		readableClassPacakgeList[i] = readableClassPackage
	}

	return readableClassPacakgeList
}

func ToReadableClassPackageList(classPackageObjectList []ClassPackage, err *CustomError) []ReadableClassPackage {
	readableClassPackageList := make([]ReadableClassPackage, len(classPackageObjectList))

	for i, item := range classPackageObjectList {
		var readableClassPackage ReadableClassPackage
		item.ToReadableClassPackage(&readableClassPackage)
		if err.IsError() {
			err.StatusCode = 500
			err.ErrorReason = "fail to parse class package"
			return nil
		}
		readableClassPackageList[i] = readableClassPackage
		
	}
	return readableClassPackageList
}
