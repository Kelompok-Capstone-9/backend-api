package models

type CreditCard struct {
	Number          string `json:"number"`
	ExpMonth        int    `json:"expire_month"`
	ExpYear         int    `json:"expire_year"`
	CVV             string `json:"cvv"`
	CreditCardToken string
}
