package database

import (
	"errors"
	"gofit-api/configs"
	"gofit-api/models"

	"gorm.io/gorm"
)

func GetMemberships(page *models.Pages, err *models.CustomError) ([]models.ReadableMembership, int) {
	var membershipObjectList []models.Membership

	result := configs.DB.Preload("User").Preload("Plan").Scopes(PaginatedQuery(page)).Find(&membershipObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableMembershipList(membershipObjectList), int(result.RowsAffected)
}

func GetMembership(membershipObject *models.Membership, err *models.CustomError) {
	result := configs.DB.Preload("User").Preload("Plan").First(membershipObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
}

func GetMembershipByUserID(userID uint, membershipObject *models.Membership, err *models.CustomError) {
	result := configs.DB.Where("user_id = ? AND is_active = true", userID).Preload("User").Preload("Plan").First(membershipObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}

}

func GetMembershipByUserName(nameUser string, err *models.CustomError) ([]models.ReadableMembership, int) {
	usersID := []int{}
	configs.DB.Model(&models.User{}).Select("id").Where("name LIKE ?", nameUser).Find(&usersID)

	var membershipObjectList []models.Membership
	result := configs.DB.Where("user_id IN ?", usersID).Preload("User").Preload("Plan").Find(&membershipObjectList)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}

	return models.ToReadableMembershipList(membershipObjectList), int(result.RowsAffected)
}

func CreateMembership(membershipObject *models.Membership, err *models.CustomError) {
	result := configs.DB.Create(membershipObject)
	if result.Error != nil {
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			err.DuplicateKey(result.Error)
		} else {
			err.FailCreateDataInDB(result.Error)
		}
	}
}

func UpdateMembership(membershipObject *models.Membership, err *models.CustomError) {
	result := configs.DB.Model(&models.Membership{}).Where("id = ?", membershipObject.ID).Updates(membershipObject)
	if result.Error != nil {
		err.FailEditDataInDB(result.Error)
	}
}

func DeleteMembership(membershipObject *models.Membership, err *models.CustomError) {
	result := configs.DB.Delete(membershipObject)
	if result.Error != nil {
		err.FailDeleteDataInDB(result.Error)
	}
}

func ActivateMembershipByID(membershipID int) error {
	err := configs.DB.Model(models.Membership{}).Where("id = ?", membershipID).Update("is_active", true).Error
	return err
}
