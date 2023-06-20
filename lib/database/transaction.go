package database

import (
	"errors"
	"gofit-api/configs"
	"gofit-api/models"

	"gorm.io/gorm"
)

func GetTransactions(offset, limit int, err *models.CustomError) ([]models.ReadableTransaction, int) {
	var transactionObjectList []models.Transaction

	result := configs.DB.Preload("User").Preload("Plan").Offset(offset).Limit(limit).Find(&transactionObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableTransactionList(transactionObjectList), int(result.RowsAffected)
}

func GetTransaction(transactionObject *models.Transaction, err *models.CustomError) {
	result := configs.DB.Preload("User").Preload("Plan").First(transactionObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
}

func GetTransactionByUserID(userID uint, transactionObject *models.Transaction, err *models.CustomError) {
	result := configs.DB.Where("user_id = ? AND is_active = true", userID).Preload("User").Preload("Plan").First(transactionObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}

}

func GetTransactionByUserName(nameUser string, err *models.CustomError) ([]models.ReadableTransaction, int) {
	usersID := []int{}
	configs.DB.Model(&models.User{}).Select("id").Where("name LIKE ?", nameUser).Find(&usersID)

	var transactionObjectList []models.Transaction
	result := configs.DB.Where("user_id IN ?", usersID).Preload("User").Preload("Plan").Find(&transactionObjectList)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}

	return models.ToReadableTransactionList(transactionObjectList), int(result.RowsAffected)
}

func CreateTransaction(transactionObject *models.Transaction, err *models.CustomError) {
	result := configs.DB.Create(transactionObject)
	if result.Error != nil {
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			err.DuplicateKey(result.Error)
		} else {
			err.FailCreateDataInDB(result.Error)
		}
	}
}

func UpdateTransaction(transactionObject *models.Transaction, err *models.CustomError) {
	result := configs.DB.Model(&models.Transaction{}).Where("id = ?", transactionObject.ID).Updates(transactionObject)
	if result.Error != nil {
		err.FailEditDataInDB(result.Error)
	}
}

func DeleteTransaction(transactionObject *models.Transaction, err *models.CustomError) {
	result := configs.DB.Delete(transactionObject)
	if result.Error != nil {
		err.FailDeleteDataInDB(result.Error)
	}
}
