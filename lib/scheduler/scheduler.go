package scheduler

import (
	"fmt"
	"gofit-api/configs"
	"gofit-api/models"

	"github.com/jasonlvhit/gocron"
)

// Check Membership Activity
func CheckMembershipActivity() {
	inactiveMembershipList := []models.Membership{}
	membershipList := []models.Membership{}
	configs.DB.Find(&membershipList)

	for _, membership := range membershipList {

		isActive := membership.CheckMembershipActivity()

		if isActive {
			continue
		} else {
			membership.IsActive = isActive
			inactiveMembershipList = append(inactiveMembershipList, membership)
		}
	}
	if len(inactiveMembershipList) != 0 {
		configs.DB.Save(&inactiveMembershipList)
	}
}

func ScheduleMembershipActivityCheck() {
	fmt.Println("Begin CrownJob")
	// gocron.Every(1).Day().Do(CheckMembershipActivity)
	// Begin job immediately upon start
	gocron.Every(1).Day().From(gocron.NextTick()).Do(CheckMembershipActivity)
	<-gocron.Start()
}
