package database

import (
	"gofit-api/configs"
	"gofit-api/models"
)

func GetLocations(offset, limit int, err *models.CustomError) ([]models.ReadableLocation, int) {
	var locationObjectList []models.Location

	result := configs.DB.Offset(offset).Limit(limit).Find(&locationObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableLocationList(locationObjectList, err), int(result.RowsAffected)
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

// func GetInstructor(instructorObject *models.Instructor, err *models.CustomError) {
// 	result := configs.DB.First(instructorObject)
// 	if result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			err.NoRecordFound(result.Error)
// 		} else {
// 			err.FailRetrieveDataFromDB(result.Error)
// 		}
// 	}
// }

// func CreateInstructor(instructorObject *models.Instructor, err *models.CustomError) {
// 	result := configs.DB.Create(instructorObject)
// 	if result.Error != nil {
// 		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
// 			err.DuplicateKey(result.Error)
// 		} else {
// 			err.FailCreateDataInDB(result.Error)
// 		}
// 	}
// }

// func UpdateInstructor(instructorObject *models.Instructor, err *models.CustomError) {
// 	result := configs.DB.Save(instructorObject)
// 	if result.Error != nil {
// 		err.FailEditDataInDB(result.Error)
// 	}
// }

// func DeleteInstructor(instructorObject *models.Instructor, err *models.CustomError) {
// 	result := configs.DB.Delete(instructorObject)
// 	if result.Error != nil {
// 		err.FailDeleteDataInDB(result.Error)
// 	}
// }
