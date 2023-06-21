package models

type LoginRequest struct{

	Email string `json:"email"`
	Password string `json:"password"`
}

type ForgotRequest struct{
	Email string `json:"email"`
	
}

type ClassTicketPaymentRequest struct{
	ReadablePaymentMethod `json:"payment_method"`
	CreditCard `json:"credit_card"`
}