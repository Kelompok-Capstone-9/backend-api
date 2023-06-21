package database

import (
	"errors"
	"gofit-api/configs"
	"gofit-api/models"

	"gorm.io/gorm"
)

func GetLocations(page *models.Pages, err *models.CustomError) ([]models.ReadableLocation, int) {
	var locationObjectList []models.Location

	result := configs.DB.Scopes(PaginatedQuery(page)).Find(&locationObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableLocationList(locationObjectList, err), int(result.RowsAffected)
}

func GetLocationsWithParams(query string, page *models.Pages, err *models.CustomError) ([]models.ReadableLocation, int) {
	var locationObjectList []models.Location

	result := configs.DB.Where(query).Scopes(PaginatedQuery(page)).Find(&locationObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableLocationList(locationObjectList, err), int(result.RowsAffected)
}

func GetIDLocationsWithParams(query string, page *models.Pages, err *models.CustomError) []int {
	var locationsIDs []int

	result := configs.DB.Model(&models.Location{}).Select("id").Where(query).Find(&locationsIDs)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil
	}

	return locationsIDs
}

func GetLocation(locationObject *models.Location, err *models.CustomError) {
	result := configs.DB.First(locationObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
}

func CreateLocation(locationObject *models.Location, err *models.CustomError) {
	result := configs.DB.Create(locationObject)
	if result.Error != nil {
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			err.DuplicateKey(result.Error)
		} else {
			err.FailCreateDataInDB(result.Error)
		}
	}
}

func UpdateLocation(locationObject *models.Location, err *models.CustomError) {
	result := configs.DB.Save(locationObject)
	if result.Error != nil {
		err.FailEditDataInDB(result.Error)
	}
}

func DeleteLocation(locationObject *models.Location, err *models.CustomError) {
	result := configs.DB.Delete(locationObject)
	if result.Error != nil {
		err.FailDeleteDataInDB(result.Error)
	}
}
