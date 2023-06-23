package models

import (
	"fmt"
	"strconv"
)

// Transaction represents the transaction object in the database
type Transaction struct {
	ID              uint          `gorm:"column:id"`
	Product         ProductType   `gorm:"type:enum('membership','class')"`
	ProductID       int           `gorm:"column:product_id"`
	Amount          int           `gorm:"column:amount"`
	TransactionCode string        `gorm:"column:transaction_code"`
	PaymentMethodID *int          `gorm:"column:payment_method_id"`
	PaymentMethod   PaymentMethod `gorm:"constraint:OnUpdate:CASCADE"`
	Status          string        `gorm:"type:enum('completed','pending', 'cancel');default:pending"`
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

func (t *Transaction) GenerateTransactionCode() {
	switch t.Product {
	case MembershipProduct:
		t.TransactionCode = fmt.Sprintf("TM%d", t.ProductID)
	case ClassProduct:
		t.TransactionCode = fmt.Sprintf("TK%d", t.ProductID)
	}
}

// ToReadableTransaction converts the transaction object to readable format
func (t *Transaction) ToReadableTransaction(readableTransaction *ReadableTransaction) {
	readableTransactionMetadata := t.Metadata.ToReadableMetadata()

	var readablePaymentMethod ReadablePaymentMethod
	t.PaymentMethod.ToReadablePaymentMethod(&readablePaymentMethod)

	readableTransaction.ID = int(t.ID)
	readableTransaction.Product = string(t.Product)
	readableTransaction.ProductID = t.ProductID
	readableTransaction.Amount = t.Amount
	readableTransaction.TransactionCode = t.TransactionCode
	readableTransaction.Status = t.Status
	readableTransaction.PaymentMethod = readablePaymentMethod
	readableTransaction.ReadableMetadata = *readableTransactionMetadata
}

// ReadableTransaction represents the readable transaction data
type ReadableTransaction struct {
	ID               int                   `json:"id"`
	Product          string                `json:"product"`
	ProductID        int                   `json:"product_id"`
	Amount           int                   `json:"amount"`
	TransactionCode  string                `json:"transaction_code"`
	PaymentMethod    ReadablePaymentMethod `json:"payment_method"`
	Status           string                `json:"status"`
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
	var transactionProductType ProductType
	transactionProductType, err.ErrorMessage = GenerateProductType(rt.Product)
	if err.IsError() {
		err.NewError(400, err.ErrorMessage, "invalid product type")
	}

	var paymentMethodObject PaymentMethod
	rt.PaymentMethod.ToPaymentMethodObject(&paymentMethodObject)
	var paymentMethodID *int = new(int)
	*paymentMethodID = int(paymentMethodObject.ID)

	transactionObject.ID = uint(rt.ID)
	transactionObject.Product = transactionProductType
	transactionObject.ProductID = rt.ProductID
	transactionObject.Amount = rt.Amount
	transactionObject.TransactionCode = rt.TransactionCode
	transactionObject.PaymentMethodID = paymentMethodID
	transactionObject.PaymentMethod = paymentMethodObject
	transactionObject.Status = rt.Status
	// transactionObject.CreatedAt, _ = time.Parse(constants.DATETIME_FORMAT, rt.CreatedAt)
	// transactionObject.UpdatedAt, _ = time.Parse(constants.DATETIME_FORMAT, rt.UpdatedAt)
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

type MidtransCallBack struct {
	TransactionStatus string `json:"transaction_status"`
	TransactionCode   string `json:"order_id"`
}
