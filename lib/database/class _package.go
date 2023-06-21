package database

import (
	"errors"
	"fmt"
	"gofit-api/configs"
	"gofit-api/models"

	"gorm.io/gorm"
)

func ClassPackageTotalData() int {
	var totalData int64
	configs.DB.Table("class_packages").Count(&totalData)
	return int(totalData)
}

func GetClassPackages(page *models.Pages, err *models.CustomError) ([]models.ReadableClassPackage, int) {
	var classpackageObjectList []models.ClassPackage

	result := configs.DB.Scopes(PaginatedQuery(page)).Preload("Class.Location").Find(&classpackageObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableClassPackageList(classpackageObjectList, err), int(result.RowsAffected)
}

func GetClassPackagesOnly(offset, limit int, err *models.CustomError) ([]models.ReadableClassPackageOnly, int) {
	var classpackageObjectList []models.ClassPackage

	result := configs.DB.Offset(offset).Limit(limit).Find(&classpackageObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableClassPackageOnlyList(classpackageObjectList, err), int(result.RowsAffected)
}

func GetClassPackagesWithParams(query string, page *models.Pages, err *models.CustomError) ([]models.ReadableClassPackage, int) {
	var classPackageObjectList []models.ClassPackage

	fmt.Println(query)

	result := configs.DB.Where(query).Scopes(PaginatedQuery(page)).Find(&classPackageObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableClassPackageList(classPackageObjectList, err), int(result.RowsAffected)
}

func GetClassPackage(classPackageObject *models.ClassPackage, err *models.CustomError) {
	result := configs.DB.Preload("Class.Location").First(classPackageObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
}

func CreateClassPackage(classPackageObject *models.ClassPackage, err *models.CustomError) {
	result := configs.DB.Create(classPackageObject)
	if result.Error != nil {
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			err.DuplicateKey(result.Error)
		} else {
			err.FailCreateDataInDB(result.Error)
		}
	}
}

func UpdateClassPackage(classPackageObject *models.ClassPackage, err *models.CustomError) {
	result := configs.DB.Save(classPackageObject)
	if result.Error != nil {
		err.FailEditDataInDB(result.Error)
	}
}

func DeleteClassPackage(classPackageObject *models.ClassPackage, err *models.CustomError) {
	result := configs.DB.Delete(classPackageObject)
	if result.Error != nil {
		err.FailDeleteDataInDB(result.Error)
	}
}
