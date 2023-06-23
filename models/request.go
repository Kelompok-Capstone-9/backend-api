package models

import "errors"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotRequest struct {
	Email string `json:"email"`
}

type PaymentRequest struct {
	ReadablePaymentMethod `json:"payment_method"`
	CreditCard            `json:"credit_card"`
}

func (pr *PaymentRequest) Validate() error {
	var err error
	switch pr.ReadablePaymentMethod.Name {
	case "gopay":
		// do something for gopay
	case "credit_card":
		switch {
		case pr.CreditCard.CVV == "":
			err = errors.New("credit card cvv cant be blank")
		case pr.CreditCard.ExpYear == 0:
			err = errors.New("credit card expire year cant be blank")
		case pr.CreditCard.ExpMonth == 0:
			err = errors.New("credit card expire month cant be blank")
		case pr.CreditCard.Number == "":
			err = errors.New("credit card number cant be blank")
		}
	case "shoope_pay":
		// do something for shoopepay
	default:
		err = errors.New("payment method name must containt correct payment method")
	}
	return err
}
