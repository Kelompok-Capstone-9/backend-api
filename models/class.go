package models

import (
	"gofit-api/constants"
	"time"
)

type Class struct {
	ID          uint
	Name        string
	Description string
	ClassType   ClassType `gorm:"type:enum('offline','online')"`
	StartedAt   time.Time
	Location    Location
	ClassPackages []ClassPackage
	Metadata    `gorm:"embedded"`
}

func (c *Class) ToReadableClass(readableClass *ReadableClass) {
	readableLocationMetadata := c.Location.ToReadableMetadata()
	readableClassMetadata := c.Metadata.ToReadableMetadata()

	readableClass.ID = int(c.ID)
	readableClass.Name = c.Name
	readableClass.Description = c.Description
	readableClass.ClassType = string(c.ClassType)
	readableClass.StartedAt = c.StartedAt.Format(constants.DATETIME_FORMAT)
	readableClass.Location.ID = int(c.Location.ID)
	readableClass.Location.Name = c.Location.Name
	readableClass.Location.Link = c.Location.Link
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
	StartedAt        string           `json:"started_at"`
	Location         ReadableLocation `json:"location"`
	ReadableMetadata `json:"metadata"`
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
	classObject.StartedAt = startedAt
	classObject.Location.ID = uint(rc.Location.ID)
	classObject.Location.Name = rc.Location.Name
	classObject.Location.Link = rc.Location.Link
	classObject.Location.Latitude = rc.Location.Latitude
	classObject.Location.Longitude = rc.Location.Longitude
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