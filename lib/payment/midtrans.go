package payment

import (
	"gofit-api/configs"
	"gofit-api/models"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func GenerateCreditCardToken(creditCard *models.CreditCard) (string, error) {
	midtrans.ClientKey = configs.AppConfig.MidtransClientKey
	resp, err := coreapi.CardToken(creditCard.Number, creditCard.ExpMonth, creditCard.ExpYear, creditCard.CVV)
	if err != nil {
		return "", err
	}
	return resp.TokenID, err
}

func PayWithCreditCard(order string, amount int64, cardToken string) coreapi.ChargeResponse {
	var c = coreapi.Client{}
	c.New(configs.AppConfig.MidtransServerKey, midtrans.Sandbox)

	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeCreditCard,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order,
			GrossAmt: amount,
		},
		CreditCard: &coreapi.CreditCardDetails{
			TokenID: cardToken,
		},
	}
	res, _ := c.ChargeTransaction(chargeReq)
	return *res
}

