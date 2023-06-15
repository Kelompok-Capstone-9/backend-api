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
	configs.DB.Where("is_active = true").Find(&membershipList)

	if len(membershipList) != 0 {
		for _, membership := range membershipList {

			isActive := membership.CheckMembershipActivity()

			if isActive {
				continue
			} else {
				membership.IsActive = isActive
				inactiveMembershipList = append(inactiveMembershipList, membership)
			}
		}
	}

	if len(inactiveMembershipList) != 0 {
		fmt.Println("found " , len(inactiveMembershipList) , " expired membership")
		configs.DB.Save(&inactiveMembershipList)
	}
}

func ScheduleMembershipActivityCheck() {
	fmt.Println("Begin CronJob")
	// gocron.Every(1).Day().Do(CheckMembershipActivity)
	// Begin job immediately upon start
	gocron.Every(1).Day().From(gocron.NextTick()).Do(CheckMembershipActivity)
	<-gocron.Start()
}
