package models

import (
	"strconv"
	"time"
)

// Membership Object for gorm
type Membership struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	PlanID    uint      `json:"plan_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Membership) InsertID(itemIDString string, err *CustomError) {
	var itemID uint64
	itemID, err.ErrorMessage = strconv.ParseUint(itemIDString, 10, 64)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id parameter"
	}
	m.ID = uint(itemID)
}

func (m *Membership) ToReadableMembership(readableMembership *ReadableMembership) {
	readableMembership.ID = int(m.ID)
	readableMembership.UserID = int(m.UserID)
	readableMembership.PlanID = int(m.PlanID)
	readableMembership.StartDate = m.StartDate
	readableMembership.EndDate = m.EndDate
	readableMembership.CreatedAt = m.CreatedAt
	readableMembership.UpdatedAt = m.UpdatedAt
}

// ReadableMembership Data or Readable data
type ReadableMembership struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PlanID    int       `json:"plan_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// convert id string to int
func (rm *ReadableMembership) InsertID(itemIDString string, err *CustomError) {
	rm.ID, err.ErrorMessage = strconv.Atoi(itemIDString)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id parameter: " + itemIDString
	}
}

func (rm *ReadableMembership) ToMembershipObject(membershipObject *Membership) {
	membershipObject.ID = uint(rm.ID)
	membershipObject.UserID = uint(rm.UserID)
	membershipObject.PlanID = uint(rm.PlanID)
	membershipObject.StartDate = rm.StartDate
	membershipObject.EndDate = rm.EndDate
	membershipObject.CreatedAt = rm.CreatedAt
	membershipObject.UpdatedAt = rm.UpdatedAt
}

func ToReadableMembershipList(membershipModelList []Membership) []ReadableMembership {
	readableMembershipList := make([]ReadableMembership, len(membershipModelList))

	for i, item := range membershipModelList {
		var readableMembership ReadableMembership
		item.ToReadableMembership(&readableMembership)
		readableMembershipList[i] = readableMembership
	}

	return readableMembershipList
}
