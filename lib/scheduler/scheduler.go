package scheduler

import (
	"fmt"
	"gofit-api/configs"
	"gofit-api/models"
	"time"

	"github.com/jasonlvhit/gocron"
)

// Check Membership Activity
func CheckMembershipActivity() {
	inactiveMembershipList := []models.Membership{}
	membershipList := []models.Membership{}
	configs.DB.Where("end_date > ?", time.Now()).Find(&membershipList)

	if len(membershipList) != 0 {
		for _, membership := range membershipList {
			membership.IsActive = false
		}
	}

	if len(inactiveMembershipList) != 0 {
		fmt.Println("found ", len(inactiveMembershipList), " expired membership")
		configs.DB.Save(&inactiveMembershipList)
		inactiveMembershipList = nil
	}
}

func ScheduleMembershipActivityCheck() {
	fmt.Println("Begin CronJob")
	// gocron.Every(1).Day().Do(CheckMembershipActivity)
	// Begin job immediately upon start
	gocron.Every(1).Day().From(gocron.NextTick()).Do(CheckMembershipActivity)
	<-gocron.Start()
}
