package models

import (
	"errors"
	"gofit-api/constants"
	"time"
)

type Class struct {
	ID            uint
	Name          string
	Description   string
	ClassType     ClassType `gorm:"type:enum('offline','online');default:offline"`
	StartedAt     time.Time
	Link          string
	LocationID    *uint
	Location      Location
	ClassPackages []ClassPackage
	Metadata      `gorm:"embedded"`
}

func (c *Class) ToReadableClass(readableClass *ReadableClass) {
	readableLocationMetadata := c.Location.ToReadableMetadata()
	readableClassMetadata := c.Metadata.ToReadableMetadata()

	readableClass.ID = int(c.ID)
	readableClass.Name = c.Name
	readableClass.Description = c.Description
	readableClass.ClassType = string(c.ClassType)
	readableClass.Link = c.Link
	readableClass.StartedAt = c.StartedAt.Format(constants.DATETIME_FORMAT)
	readableClass.Location.ID = int(c.Location.ID)
	readableClass.Location.Name = c.Location.Name
	readableClass.Location.City = c.Location.City
	readableClass.Location.Latitude = c.Location.Latitude
	readableClass.Location.Longitude = c.Location.Longitude
	readableClass.Location.ReadableMetadata = *readableLocationMetadata
	readableClass.ReadableMetadata = *readableClassMetadata
}

type ReadableClass struct {
	ID               int              `json:"id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	ClassType        string           `json:"class_type"`
	Link             string           `json:"link"`
	StartedAt        string           `json:"started_at"`
	Location         ReadableLocation `json:"location"`
	ReadableMetadata `json:"metadata"`
}

func (rc *ReadableClass) Validate() error {
	switch {
	case rc.Name == "":
		return errors.New("invalid name")
	case rc.StartedAt == "":
		return errors.New("invalid started at")
	}

	switch rc.ClassType {
	case "online":
		if rc.Link == "" {
			return errors.New("link for online class cant be blank")
		}
	case "offline":
		if rc.Location.ID == 0 {
			return errors.New("location for offline class cant be blank")
		}
	default:
		return errors.New("invalid class type. must containt offline or online")
	}

	return nil
}

func (rc *ReadableClass) EditValidate() error {
	allFieldBlank := rc.Name == "" && rc.ClassType == "" && rc.StartedAt == "" && rc.Location.ID == 0
	if allFieldBlank {
		return errors.New("all field is blank. nothing to change")
	}

	if rc.ClassType != "" {
		if rc.ClassType != "offline" && rc.ClassType != "online" {
			return errors.New("invalid class type. must containt offline or online")
		}
	}

	return nil
}

func (rc *ReadableClass) ToClassObject(classObject *Class, err *CustomError) {
	var classType ClassType
	classType, err.ErrorMessage = GenerateClassType(rc.ClassType)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid class type input"
	}

	var startedAt time.Time
	startedAt, err.ErrorMessage = time.Parse(constants.DATETIME_FORMAT, rc.StartedAt)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "fail to parse started at time"
	}

	classObject.ID = uint(rc.ID)
	classObject.Name = rc.Name
	classObject.Description = rc.Description
	classObject.ClassType = classType
	classObject.Link = rc.Link
	classObject.StartedAt = startedAt
	classObject.Location.ID = uint(rc.Location.ID)
}

func ToReadableClassList(classObjectList []Class, err *CustomError) []ReadableClass {
	readableClassList := make([]ReadableClass, len(classObjectList))

	for i, item := range classObjectList {
		var readableClass ReadableClass
		item.ToReadableClass(&readableClass)
		if err.IsError() {
			err.StatusCode = 500
			err.ErrorReason = "fail to parse location"
			return nil
		}
		readableClassList[i] = readableClass
	}

	return readableClassList
}
