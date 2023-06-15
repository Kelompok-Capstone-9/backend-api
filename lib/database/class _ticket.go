package database

import (
	"errors"
	"gofit-api/configs"
	"gofit-api/models"

	"gorm.io/gorm"
)

func GetClassTickets(offset, limit int, err *models.CustomError) ([]models.ReadableClassTicket, int) {
	var classTicketObjectList []models.ClassTicket

	result := configs.DB.Offset(offset).Limit(limit).Preload("User").Preload("ClassPackage.Class.Location").Find(&classTicketObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableClassTicketList(classTicketObjectList, err), int(result.RowsAffected)
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

func GetClassTicket(classTicketObject *models.ClassTicket, err *models.CustomError) {
	result := configs.DB.Preload("User").Preload("ClassPackage.Class.Location").First(classTicketObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
}

func CreateClassTicket(classTicketObject *models.ClassTicket, err *models.CustomError) {
	result := configs.DB.Create(classTicketObject)
	if result.Error != nil {
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			err.DuplicateKey(result.Error)
		} else {
			err.FailCreateDataInDB(result.Error)
		}
	}
}

func UpdateClassTicket(classTicketObject *models.ClassTicket, err *models.CustomError) {
	result := configs.DB.Save(classTicketObject)
	if result.Error != nil {
		err.FailEditDataInDB(result.Error)
	}
}

func DeleteClassTicket(classTicketObject *models.ClassTicket, err *models.CustomError) {
	result := configs.DB.Delete(classTicketObject)
	if result.Error != nil {
		err.FailDeleteDataInDB(result.Error)
	}
}
