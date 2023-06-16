package models

import (
	"errors"
)

type Instructor struct {
	ID          uint
	Name        string
	Description string
	Metadata    `gorm:"embedded"`
}

func (i *Instructor) ToReadableInstructor(readableInstructor *ReadableInstructor) {
	readableMetadata := i.ToReadableMetadata()

	readableInstructor.ID = int(i.ID)
	readableInstructor.Name = i.Name
	readableInstructor.Description = i.Description
	readableInstructor.ReadableMetadata = *readableMetadata
}

type ReadableInstructor struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ReadableMetadata `json:"metadata"`
}

func (ri *ReadableInstructor) ToInstructorObject(instructorObject *Instructor) {
	instructorObject.ID = uint(ri.ID)
	instructorObject.Name = ri.Name
	instructorObject.Description = ri.Description
}

func (ri *ReadableInstructor) Validate() error {
	switch {
	case ri.Name == "":
		return errors.New("invalid name")
	case ri.Description == "":
		return errors.New("invalid description")
	}
	return nil
}

func (ri *ReadableInstructor) EditValidate() error {
	if ri.Name == "" && ri.Description == "" {
		return errors.New("invalid name and description field")
	}
	return nil
}

func ToReadableInstructorList(instructorObjectList []Instructor, err *CustomError) []ReadableInstructor {
	readableInstructorList := make([]ReadableInstructor, len(instructorObjectList))

	for i, item := range instructorObjectList {
		var readableInstructor ReadableInstructor
		item.ToReadableInstructor(&readableInstructor)
		if err.IsError() {
			err.StatusCode = 500
			err.ErrorReason = "fail to parse instructor"
			return nil
		}
		readableInstructorList[i] = readableInstructor
	}

	return readableInstructorList
}
