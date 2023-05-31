package models

import (
	"gofit-api/constants"
	"strconv"
)

// Plan struct for gorm
type Plan struct {
	ID          uint `gorm:"column:id"`
	Name        string
	Duration    int
	Price       int
	Memberships []Membership `gorm:"constraint:OnUpdate:CASCADE"`
	Metadata    `gorm:"embedded"`
}

// ReadablePlan represents readable plan data
type ReadablePlan struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Duration         int    `json:"duration"`
	Price            int    `json:"price"`
	ReadableMetadata `json:"metadata"`
}

func (u *Plan) InsertID(itemIDString string, err *CustomError) {
	var itemID int
	itemID, err.ErrorMessage = strconv.Atoi(itemIDString)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id paramater"
	}
	u.ID = uint(itemID)
}

// Convert ID string to int
func (rp *ReadablePlan) InsertID(itemIDString string, err *CustomError) {
	rp.ID, err.ErrorMessage = strconv.Atoi(itemIDString)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id parameter: " + itemIDString
	}
}

func (rp *ReadablePlan) ToReadablePlan(planObject *Plan) {
	rp.ID = int(planObject.ID)
	rp.Name = planObject.Name
	rp.Duration = planObject.Duration
	rp.Price = planObject.Price
	rp.ReadableMetadata.CreatedAt = planObject.Metadata.CreatedAt.Format(constants.DATETIME_FORMAT)
	rp.ReadableMetadata.UpdatedAt = planObject.Metadata.UpdatedAt.Format(constants.DATETIME_FORMAT)
}

// ToReadablePlanList converts a list of Plan models to a list of ReadablePlan models
func ToReadablePlanList(planModelList []Plan) []ReadablePlan {
	readablePlanList := make([]ReadablePlan, len(planModelList))

	for i, item := range planModelList {
		var readablePlan ReadablePlan
		readablePlan.ID = int(item.ID)
		readablePlan.Name = item.Name
		readablePlan.Duration = item.Duration
		readablePlan.Price = item.Price
		// readablePlan.ReadableMetadata = *metadata
		readablePlanList[i] = readablePlan
	}

	return readablePlanList
}
