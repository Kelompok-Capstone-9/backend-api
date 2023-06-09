package database

import (
	"errors"
	"gofit-api/configs"
	"gofit-api/models"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var mysqlErr *mysql.MySQLError

func GetUsers(offset, limit int, err *models.CustomError) ([]models.ReadableUser, int) {
	var userObjectList []models.User

	result := configs.DB.Offset(offset).Limit(limit).Find(&userObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableUserList(userObjectList, err), int(result.RowsAffected)
}

func GetUser(userObject *models.User, err *models.CustomError) {
	result := configs.DB.First(userObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
}

func CreateUser(userObject *models.User, err *models.CustomError) {
	userObject.HashingPassword(err)
	result := configs.DB.Create(userObject)
	if result.Error != nil {
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			err.DuplicateKey(result.Error)
		} else {
			err.FailCreateDataInDB(result.Error)
		}
	}
}

func UpdateUser(userObject *models.User, err *models.CustomError) {
	result := configs.DB.Save(userObject)
	if result.Error != nil {
		err.FailEditDataInDB(result.Error)
	}
}

func DeleteUser(userObject *models.User, err *models.CustomError) {
	result := configs.DB.Delete(userObject)
	if result.Error != nil {
		err.FailDeleteDataInDB(result.Error)
	}
}

func Login(email string, err *models.CustomError) models.User {
	var userObject models.User
	result := configs.DB.Where("email = ?", email).First(&userObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.FailLogin()
			return models.User{}
		} else {
			err.ErrorMessage = result.Error
			err.FailRetrieveDataFromDB(result.Error)
			return models.User{}
		}
	}
	return userObject
}
