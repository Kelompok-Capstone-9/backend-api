package models

import (
	"errors"
	"strconv"
)

// PaymentMethod struct for gorm
type PaymentMethod struct {
	ID           uint `gorm:"column:id"`
	Name         string
	Transactions []Transaction `gorm:"constraint:OnUpdate:CASCADE"`
	Metadata     `gorm:"embedded"`
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

func (pm *PaymentMethod) ToReadablePaymentMethod(readablePaymentMethod *ReadablePaymentMethod) {
	readableMetadata := pm.ToReadableMetadata()
	readablePaymentMethod.ID = int(pm.ID)
	readablePaymentMethod.Name = pm.Name
	readablePaymentMethod.ReadableMetadata = *readableMetadata

}

// ReadablePaymentMethod represents readable payment method data
type ReadablePaymentMethod struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	ReadableMetadata `json:"metadata"`
}

func (rpm *ReadablePaymentMethod) Validate() error {
	if rpm.Name == "" {
		return errors.New("name field is blank")
	}
	return nil
}

func (rpm *ReadablePaymentMethod) ToPaymentMethodObject(paymentObject *PaymentMethod) {
	paymentObject.ID = uint(rpm.ID)
	paymentObject.Name = rpm.Name
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
