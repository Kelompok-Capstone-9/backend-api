package database

import (
	"errors"
	"gofit-api/configs"
	"gofit-api/models"

	"gorm.io/gorm"
)

func GetHealthtips(offset, limit int, err *models.CustomError) ([]models.ReadableHealthtip, int) {
	var healthtipObjectList []models.Healthtip

	result := configs.DB.Offset(offset).Limit(limit).Preload("User").Find(&healthtipObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableHealthtipList(healthtipObjectList, err), int(result.RowsAffected)
}

func GetHealthtip(healthtipObject *models.Healthtip, err *models.CustomError) {
	result := configs.DB.Preload("User").First(healthtipObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
}

func CreateHealthtip(healthtipObject *models.Healthtip, err *models.CustomError) {
	result := configs.DB.Create(healthtipObject)
	if result.Error != nil {
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			err.DuplicateKey(result.Error)
		} else {
			err.FailCreateDataInDB(result.Error)
		}
	}
}

func UpdateHealthtip(healthtipObject *models.Healthtip, err *models.CustomError) {
	result := configs.DB.Save(healthtipObject)
	if result.Error != nil {
		err.FailEditDataInDB(result.Error)
	}
}

func DeleteHealthtip(healthtipObject *models.Healthtip, err *models.CustomError) {
	result := configs.DB.Delete(healthtipObject)
	if result.Error != nil {
		err.FailDeleteDataInDB(result.Error)
	}
}
