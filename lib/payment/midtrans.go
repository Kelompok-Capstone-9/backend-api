package payment

import (
	"fmt"
	"gofit-api/configs"
	"gofit-api/models"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

func GenerateCreditCardToken(creditCard *models.CreditCard) (string, error) {
	midtrans.ClientKey = configs.AppConfig.MidtransClientKey
	resp, err := coreapi.CardToken(creditCard.Number, creditCard.ExpMonth, creditCard.ExpYear, creditCard.CVV)
	if err != nil {
		return "", err.RawError
	}
	return resp.TokenID, nil
}

func PayWithCreditCard(chargeRequest models.ChargeAPIRequest) (coreapi.ChargeResponse, *midtrans.Error) {
	var c = coreapi.Client{}
	c.New(configs.AppConfig.MidtransServerKey, midtrans.Sandbox)

	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeCreditCard,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  chargeRequest.OrderID,
			GrossAmt: chargeRequest.GrossAmount,
		},
		CreditCard: &coreapi.CreditCardDetails{
			TokenID: chargeRequest.CreditCardToken,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    chargeRequest.ItemDetails[0].ID,
				Price: chargeRequest.ItemDetails[0].Price,
				Qty:   1,
				Name:  chargeRequest.ItemDetails[0].Name,
			},
		},
	}

	res, err := c.ChargeTransaction(chargeReq)
	if err != nil {
		fmt.Println(err)
		return coreapi.ChargeResponse{}, err
	}
	return *res, nil
}
