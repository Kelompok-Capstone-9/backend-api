package database

import (
	"errors"
	"gofit-api/configs"
	"gofit-api/models"
	"log"

	"gorm.io/gorm"
)

func GetPlans(offset, limit int, err *models.CustomError) ([]models.ReadablePlan, int) {
	var planObjectList []models.Plan

	result := configs.DB.Offset(offset).Limit(limit).Find(&planObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadablePlanList(planObjectList), int(result.RowsAffected)
}

func GetPlan(planObject *models.Plan, err *models.CustomError) {
	result := configs.DB.First(planObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
}

func CreatePlan(planObject *models.Plan, err *models.CustomError) {
	result := configs.DB.Create(planObject)
	if result.Error != nil {
		log.Println("Error creating plan:", result.Error)
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			err.DuplicateKey(result.Error)
		} else {
			err.FailCreateDataInDB(result.Error)
		}
	}
}

func UpdatePlan(planObject *models.Plan, err *models.CustomError) {
	result := configs.DB.Save(planObject)
	if result.Error != nil {
		err.FailEditDataInDB(result.Error)
	}
}

func DeletePlan(planObject *models.Plan, err *models.CustomError) {
	result := configs.DB.Delete(planObject)
	if result.Error != nil {
		err.FailDeleteDataInDB(result.Error)
	}
}
