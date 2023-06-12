package database

import (
	"errors"
	"gofit-api/configs"
	"gofit-api/models"

	"gorm.io/gorm"
)

func GetClassPackages(offset, limit int, err *models.CustomError) ([]models.ReadableClassPackage, int) {
	var classpackageObjectList []models.ClassPackage

	result := configs.DB.Offset(offset).Limit(limit).Preload("Class.Location").Find(&classpackageObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableClassPackageList(classpackageObjectList, err), int(result.RowsAffected)
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

func GetClassPackage(classObject *models.Class, err *models.CustomError) {
	result := configs.DB.Preload("Location").First(classObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
}

func CreateClassPackage(classObject *models.Class, err *models.CustomError) {
	result := configs.DB.Create(classObject)
	if result.Error != nil {
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			err.DuplicateKey(result.Error)
		} else {
			err.FailCreateDataInDB(result.Error)
		}
	}
}

func UpdateClassPackage(classObject *models.Class, err *models.CustomError) {
	result := configs.DB.Save(classObject)
	if result.Error != nil {
		err.FailEditDataInDB(result.Error)
	}
}

func DeleteClassPackage(classObject *models.Class, err *models.CustomError) {
	result := configs.DB.Delete(classObject)
	if result.Error != nil {
		err.FailDeleteDataInDB(result.Error)
	}
}
