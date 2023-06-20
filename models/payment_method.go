package models

import (
	"gofit-api/constants"
	"strconv"
)

// PaymentMethod struct for gorm
type PaymentMethod struct {
	ID           uint `gorm:"column:id"`
	Name         string
	Transactions []Transaction `gorm:"constraint:OnUpdate:CASCADE"`
	Metadata     `gorm:"embedded"`
}

// ReadablePaymentMethod represents readable payment method data
type ReadablePaymentMethod struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	ReadableMetadata `json:"metadata"`
}

func (u *PaymentMethod) InsertID(itemIDString string, err *CustomError) {
	var itemID int
	itemID, err.ErrorMessage = strconv.Atoi(itemIDString)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id paramater"
	}
	u.ID = uint(itemID)
}

// Convert ID string to int
func (rpm *ReadablePaymentMethod) InsertID(itemIDString string, err *CustomError) {
	rpm.ID, err.ErrorMessage = strconv.Atoi(itemIDString)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id parameter: " + itemIDString
	}
}

func (rpm *ReadablePaymentMethod) ToReadablePaymentMethod(paymentMethodObject *PaymentMethod) {
	rpm.ID = int(paymentMethodObject.ID)
	rpm.Name = paymentMethodObject.Name
	rpm.ReadableMetadata.CreatedAt = paymentMethodObject.Metadata.CreatedAt.Format(constants.DATETIME_FORMAT)
	rpm.ReadableMetadata.UpdatedAt = paymentMethodObject.Metadata.UpdatedAt.Format(constants.DATETIME_FORMAT)
}

// ToReadablePaymentMethodList converts a list of PaymentMethod models to a list of ReadablePaymentMethod models
func ToReadablePaymentMethodList(paymentMethodModelList []PaymentMethod) []ReadablePaymentMethod {
	readablePaymentMethodList := make([]ReadablePaymentMethod, len(paymentMethodModelList))

	for i, item := range paymentMethodModelList {
		var readablePaymentMethod ReadablePaymentMethod
		readablePaymentMethod.ID = int(item.ID)
		readablePaymentMethod.Name = item.Name
		// readablePaymentMethod.ReadableMetadata = *metadata
		readablePaymentMethodList[i] = readablePaymentMethod
	}

	return readablePaymentMethodList
}
