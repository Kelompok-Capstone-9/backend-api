package database

import (
	"errors"
	"gofit-api/configs"
	"gofit-api/models"

	"gorm.io/gorm"
)

func GetMemberships(offset, limit int, err *models.CustomError) ([]models.ReadableMembership, int) {
	var membershipObjectList []models.Membership

	result := configs.DB.Offset(offset).Limit(limit).Find(&membershipObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableMembershipList(membershipObjectList), int(result.RowsAffected)
}

func GetMembership(membershipObject *models.Membership, err *models.CustomError) {
	result := configs.DB.First(membershipObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
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
	result := configs.DB.Save(membershipObject)
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
