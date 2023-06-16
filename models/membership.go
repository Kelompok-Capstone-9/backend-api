package models

import (
	"gofit-api/constants"
	"strconv"
	"time"
)

// Membership Object for gorm
type Membership struct {
	ID        uint      `gorm:"column:id"`
	UserID    uint      `gorm:"column:user_id"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE"`
	PlanID    uint      `gorm:"column:plan_id"`
	Plan      Plan      `gorm:"constraint:OnUpdate:CASCADE"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool
	Metadata  `gorm:"embedded"`
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
	readableUserMetadata := m.User.ToReadableMetadata()
	readableMembershipMetadata := m.ToReadableMetadata()
	readablePlanMetadata := m.Plan.ToReadableMetadata()
	readableMembership.ID = int(m.ID)
	readableMembership.User.ID = int(m.User.ID)
	readableMembership.User.Name = m.User.Name
	readableMembership.User.Email = m.User.Email
	readableMembership.User.Password = "********"
	readableMembership.User.Gender = string(m.User.Gender)
	readableMembership.User.Height = m.User.Height
	readableMembership.User.GoalHeight = m.User.GoalHeight
	readableMembership.User.Weight = m.User.Weight
	readableMembership.User.GoalWeight = m.User.GoalWeight
	readableMembership.User.ReadableMetadata = *readableUserMetadata
	readableMembership.Plan.ID = int(m.Plan.ID)
	readableMembership.Plan.Name = m.Plan.Name
	readableMembership.Plan.Duration = m.Plan.Duration
	readableMembership.Plan.Price = m.Plan.Price
	readableMembership.Plan.Description = m.Plan.Description
	readableMembership.Plan.ReadableMetadata = *readableMembershipMetadata
	readableMembership.StartDate = m.StartDate.Format(constants.DATETIME_FORMAT)
	readableMembership.EndDate = m.EndDate.Format(constants.DATETIME_FORMAT)
	readableMembership.IsActive = m.CheckMembershipActivity()
	readableMembership.ReadableMetadata = *readablePlanMetadata
}

// ReadableMembership Data or Readable data
type ReadableMembership struct {
	ID               int          `json:"id"`
	User             ReadableUser `json:"user"`
	Plan             ReadablePlan `json:"plan"`
	StartDate        string       `json:"start_date"`
	EndDate          string       `json:"end_date"`
	IsActive         bool         `json:"is_active"`
	ReadableMetadata `json:"metadata"`
}

// convert id string to int
func (rm *ReadableMembership) InsertID(itemIDString string, err *CustomError) {
	rm.ID, err.ErrorMessage = strconv.Atoi(itemIDString)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id parameter: " + itemIDString
	}
}

// IsActive checks if the membership is active based on current date
func (m *Membership) CheckMembershipActivity() bool {
	currentTime := time.Now()
	return m.StartDate.After(currentTime) && currentTime.Before(m.EndDate)
}

func (rm *ReadableMembership) ToMembershipObject(membershipObject *Membership, err *CustomError) {
	membershipObject.ID = uint(rm.ID)
	membershipObject.UserID = uint(rm.User.ID)
	membershipObject.PlanID = uint(rm.Plan.ID)
	var started, ended time.Time
	started, err.ErrorMessage = time.Parse(constants.DATETIME_FORMAT, rm.StartDate)
	if err.ErrorMessage != nil {
		err.StatusCode = 400
		err.ErrorReason = "fail to parse metadata time"
	}
	ended, err.ErrorMessage = time.Parse(constants.DATETIME_FORMAT, rm.EndDate)
	if err.ErrorMessage != nil {
		err.StatusCode = 400
		err.ErrorReason = "fail to parse metadata time"
	}
	membershipObject.StartDate = started
	membershipObject.EndDate = ended
	// membershipObject.CreatedAt = rm.CreatedAt
	// membershipObject.UpdatedAt = rm.UpdatedAt
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
