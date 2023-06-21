package database

import (
	"errors"
	"gofit-api/configs"
	"gofit-api/models"

	"gorm.io/gorm"
)

func ClassTotalData() int {
	var totalData int64
	configs.DB.Table("classes").Count(&totalData)
	return int(totalData)
}

func GetClasses(page *models.Pages, err *models.CustomError) ([]models.ReadableClass, int) {
	var classObjectList []models.Class

	result := configs.DB.Scopes(PaginatedQuery(page)).Preload("ClassPackages").Preload("Location").Find(&classObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableClassList(classObjectList, err), int(result.RowsAffected)
}

// func GetInstructorsWithParams(params *models.GeneralParameter, err *models.CustomError) ([]models.ReadableInstructor, int) {
// 	var instructorObjectList []models.Instructor

// 	result := configs.DB.Where("name LIKE ?", params.Name).Offset(params.Page.Offset).Limit(params.Page.Limit).Find(&instructorObjectList)
// 	if result.Error != nil {
// 		err.FailRetrieveDataFromDB(result.Error)
// 		return nil, 0
// 	}

// 	return models.ToReadableInstructorList(instructorObjectList, err), int(result.RowsAffected)
// }

func GetClass(classObject *models.Class, err *models.CustomError) {
	result := configs.DB.Preload("ClassPackages").Preload("Location").First(classObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
}

func CreateClass(classObject *models.Class, err *models.CustomError) {
	result := configs.DB.Create(classObject)
	if result.Error != nil {
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			err.DuplicateKey(result.Error)
		} else {
			err.FailCreateDataInDB(result.Error)
		}
	}
}

func UpdateClass(classObject *models.Class, err *models.CustomError) {
	result := configs.DB.Save(classObject)
	if result.Error != nil {
		err.FailEditDataInDB(result.Error)
	}
}

func DeleteClass(classObject *models.Class, err *models.CustomError) {
	result := configs.DB.Delete(classObject)
	if result.Error != nil {
		err.FailDeleteDataInDB(result.Error)
	}
}
