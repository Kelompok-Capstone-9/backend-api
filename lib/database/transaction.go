package database

import (
	"errors"
	"gofit-api/configs"
	"gofit-api/models"

	"gorm.io/gorm"
)

func GetTransactions(page *models.Pages, err *models.CustomError) ([]models.ReadableTransaction, int) {
	var transactionObjectList []models.Transaction

	result := configs.DB.Scopes(PaginatedQuery(page)).Preload("PaymentMethod").Find(&transactionObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadableTransactionList(transactionObjectList), int(result.RowsAffected)
}

func GetTransaction(transactionObject *models.Transaction, err *models.CustomError) {
	result := configs.DB.Preload("PaymentMethod").First(transactionObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
}

func GetTransactionByUserID(userID uint, transactionObject *models.Transaction, err *models.CustomError) {
	result := configs.DB.Preload("PaymentMethod").First(transactionObject)
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
	configs.DB.Model(&models.User{}).Select("id").Where("name LIKE ?", nameUser).Preload("PaymentMethod").Find(&usersID)

	var transactionObjectList []models.Transaction
	result := configs.DB.Find(&transactionObjectList)
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

func GetTransactionByCode(transactionCode string, err *models.CustomError) models.Transaction {
	var transaction models.Transaction
	result := configs.DB.Where("transaction_code = ?", transactionCode).First(&transaction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
	return transaction
}
