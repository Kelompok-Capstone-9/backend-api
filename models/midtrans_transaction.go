package models

type TransactionDetail struct {
	OrderID     string `json:"order_id"`
	GrossAmount int64  `json:"gross_amount"`
}

type ItemDetail struct {
	ID       string `json:"id"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

type CustomerDetails struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
}

type ChargeAPIRequest struct {
	PaymentType       string `json:"payment_type"`
	TransactionDetail `json:"transaction_details"`
	CreditCard        `json:"credit_card"`
	CustomerDetails   `json:"customer_details"`
	ItemDetails       []ItemDetail `json:"item_details"`
}
