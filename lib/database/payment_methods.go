package database

import (
	"errors"
	"gofit-api/configs"
	"gofit-api/models"
	"log"

	"gorm.io/gorm"
)

func GetPaymentMethods(offset, limit int, err *models.CustomError) ([]models.ReadablePaymentMethod, int) {
	var paymentMethodObjectList []models.PaymentMethod

	result := configs.DB.Offset(offset).Limit(limit).Find(&paymentMethodObjectList)
	if result.Error != nil {
		err.FailRetrieveDataFromDB(result.Error)
		return nil, 0
	}

	return models.ToReadablePaymentMethodList(paymentMethodObjectList), int(result.RowsAffected)
}

func GetPaymentMethod(paymentMethodObject *models.PaymentMethod, err *models.CustomError) {
	result := configs.DB.First(paymentMethodObject)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err.NoRecordFound(result.Error)
		} else {
			err.FailRetrieveDataFromDB(result.Error)
		}
	}
}

func CreatePaymentMethod(paymentMethodObject *models.PaymentMethod, err *models.CustomError) {
	result := configs.DB.Create(paymentMethodObject)
	if result.Error != nil {
		log.Println("Error creating payment method:", result.Error)
		if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
			err.DuplicateKey(result.Error)
		} else {
			err.FailCreateDataInDB(result.Error)
		}
	}
}

func UpdatePaymentMethod(paymentMethodObject *models.PaymentMethod, err *models.CustomError) {
	result := configs.DB.Save(paymentMethodObject)
	if result.Error != nil {
		err.FailEditDataInDB(result.Error)
	}
}

func DeletePaymentMethod(paymentMethodObject *models.PaymentMethod, err *models.CustomError) {
	result := configs.DB.Delete(paymentMethodObject)
	if result.Error != nil {
		err.FailDeleteDataInDB(result.Error)
	}
}
