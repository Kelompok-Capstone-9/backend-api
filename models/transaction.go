package models

import (
	"gofit-api/constants"
	"strconv"
	"time"
)

// Transaction Object for gorm
type Transaction struct {
	ID              uint   `gorm:"column:id"`
	Product         string `gorm:"type:enum('pending','completed','canceled')"`
	ProductID       int    `gorm:"column:product_id"`
	Amount          int
	PaymentMethodID int           `gorm:"column:payment_method_id"`
	PaymentMethod   PaymentMethod `gorm:"constraint:OnUpdate:CASCADE"`
	Status          string        `gorm:"type:enum('pending','completed','canceled')"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Metadata        `gorm:"embedded"`
}

func (t *Transaction) InsertID(itemIDString string, err *CustomError) {
	var itemID uint64
	itemID, err.ErrorMessage = strconv.ParseUint(itemIDString, 10, 64)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id parameter"
	}
	t.ID = uint(itemID)
}

func (t *Transaction) ToReadableTransaction(readableTransaction *ReadableTransaction) {
	readableTransaction.ID = int(t.ID)
	readableTransaction.PaymentMethod.ID = t.PaymentMethodID
	readableTransaction.PaymentMethod.Name = t.PaymentMethod.Name
	readableTransaction.ReadableMetadata.CreatedAt = t.CreatedAt.Format(constants.DATETIME_FORMAT)
	readableTransaction.ReadableMetadata.UpdatedAt = t.UpdatedAt.Format(constants.DATETIME_FORMAT)
}

// ReadableTransaction Data or Readable data
type ReadableTransaction struct {
	ID               int                   `json:"id"`
	PaymentMethod    ReadablePaymentMethod `json:"payment_method"`
	ReadableMetadata `json:"metadata"`
}

// Convert id string to int
func (rt *ReadableTransaction) InsertID(itemIDString string, err *CustomError) {
	rt.ID, err.ErrorMessage = strconv.Atoi(itemIDString)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid id parameter: " + itemIDString
	}
}

func (rm *ReadableTransaction) ToTransactionObject(transactionObject *Transaction, err *CustomError) {
	transactionObject.ID = uint(rm.ID)
	transactionObject.PaymentMethodID = rm.PaymentMethod.ID
}

func ToReadableTransactionList(transactionModelList []Transaction) []ReadableTransaction {
	readableTransactionList := make([]ReadableTransaction, len(transactionModelList))

	for i, item := range transactionModelList {
		var readableTransaction ReadableTransaction
		item.ToReadableTransaction(&readableTransaction)
		readableTransactionList[i] = readableTransaction
	}

	return readableTransactionList
}
