package models

import (
	"gofit-api/constants"

	"time"

	"gorm.io/gorm"
)

type Metadata struct {
	DeletedAt gorm.DeletedAt
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Metadata) ToReadableMetadata() *ReadableMetadata {
	return &ReadableMetadata{
		CreatedAt: m.CreatedAt.Format(constants.DATETIME_FORMAT),
		UpdatedAt: m.CreatedAt.Format(constants.DATETIME_FORMAT),
	}
}

type ReadableMetadata struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (rm *ReadableMetadata) ToMetadata(err *CustomError) *Metadata {
	var created, updated time.Time

	created, err.ErrorMessage = time.Parse(constants.DATETIME_FORMAT, rm.CreatedAt)
	if err.ErrorMessage != nil {
		err.StatusCode = 400
		err.ErrorReason = "fail to parse metadata time"
		return nil
	}
	updated, err.ErrorMessage = time.Parse(constants.DATETIME_FORMAT, rm.UpdatedAt)
	if err.ErrorMessage != nil {
		err.StatusCode = 400
		err.ErrorReason = "fail to parse metadata time"
		return nil
	}

	return &Metadata{
		CreatedAt: created,
		UpdatedAt: updated,
	}
}
