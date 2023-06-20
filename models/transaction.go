package models

import (
	"gofit-api/constants"
	"strconv"
	"time"
)

// Transaction represents the transaction object in the database
type Transaction struct {
	ID              uint          `gorm:"column:id"`
	Product         string        `gorm:"type:enum('pending','completed', 'canceled' )"`
	ProductID       int           `gorm:"column:product_id"`
	Amount          int           `gorm:"column:amount"`
	InvoiceID       string        `gorm:"column:invoice_id"`
	PaymentMethodID int           `gorm:"column:payment_method_id"`
	PaymentMethod   PaymentMethod `gorm:"constraint:OnUpdate:CASCADE"`
	Status          string        `gorm:"type:enum('pending','completed', 'canceled' )"`
	CreatedAt       time.Time     `gorm:"column:created_at"`
	UpdatedAt       time.Time     `gorm:"column:updated_at"`
	Metadata        `gorm:"embedded"`
}

// InsertID converts and inserts the item ID from string to uint
func (t *Transaction) InsertID(itemIDString string, err *CustomError) {
	var itemID uint64
	itemID, err.ErrorMessage = strconv.ParseUint(itemIDString, 10, 64)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid ID parameter"
	}
	t.ID = uint(itemID)
}

// ToReadableTransaction converts the transaction object to readable format
func (t *Transaction) ToReadableTransaction(readableTransaction *ReadableTransaction) {
	readableTransactionMetadata := t.ToReadableMetadata()
	// readablePaymentMethodMetadata := t.PaymentMethod.ToReadableMetadata()
	readableTransaction.ID = int(t.ID)
	readableTransaction.Product = t.Product
	readableTransaction.ProductID = t.ProductID
	readableTransaction.Amount = t.Amount
	readableTransaction.InvoiceID = t.InvoiceID
	readableTransaction.Status = t.Status
	readableTransaction.PaymentMethod.ReadableMetadata = *readableTransactionMetadata
}

// ReadableTransaction represents the readable transaction data
type ReadableTransaction struct {
	ID               int                   `json:"id"`
	Product          string                `json:"product"`
	ProductID        int                   `json:"product_id"`
	Amount           int                   `json:"amount"`
	InvoiceID        string                `json:"invoice_id"`
	PaymentMethod    ReadablePaymentMethod `json:"payment_method"`
	Status           string                `json:"status"`
	CreatedAt        string                `json:"created_at"`
	UpdatedAt        string                `json:"updated_at"`
	ReadableMetadata `json:"metadata"`
}

// InsertID converts and inserts the item ID from string to int
func (rt *ReadableTransaction) InsertID(itemIDString string, err *CustomError) {
	rt.ID, err.ErrorMessage = strconv.Atoi(itemIDString)
	if err.IsError() {
		err.StatusCode = 400
		err.ErrorReason = "invalid ID parameter: " + itemIDString
	}
}

// ToTransactionObject converts the readable transaction to a transaction object
func (rt *ReadableTransaction) ToTransactionObject(transactionObject *Transaction, err *CustomError) {
	transactionObject.ID = uint(rt.ID)
	transactionObject.Product = rt.Product
	transactionObject.ProductID = rt.ProductID
	transactionObject.Amount = rt.Amount
	transactionObject.InvoiceID = rt.InvoiceID
	transactionObject.Status = rt.Status
	transactionObject.CreatedAt, _ = time.Parse(constants.DATETIME_FORMAT, rt.CreatedAt)
	transactionObject.UpdatedAt, _ = time.Parse(constants.DATETIME_FORMAT, rt.UpdatedAt)
}

// ToReadableTransactionList converts a list of transaction models to a list of readable transaction models
func ToReadableTransactionList(transactionModelList []Transaction) []ReadableTransaction {
	readableTransactionList := make([]ReadableTransaction, len(transactionModelList))

	for i, item := range transactionModelList {
		var readableTransaction ReadableTransaction
		item.ToReadableTransaction(&readableTransaction)
		readableTransactionList[i] = readableTransaction
	}

	return readableTransactionList
}
