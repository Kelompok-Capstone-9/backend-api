package models

import (
	"strconv"
)

// Plan struct for gorm
type Plan struct {
	ID       uint
	Name     string
	Duration int
	Price    int
	Metadata `gorm:"embedded"`
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

func (rp *ReadablePlan) ToPlanObject(planObject *Plan) {
	planObject.ID = uint(rp.ID)
	planObject.Name = rp.Name
	planObject.Duration = rp.Duration
	planObject.Price = rp.Price
	// planObject.Metadata = *metadata
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
