package service

import (
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81"
)



type StripeClient struct {
	SecretKey string
}

func NewStripeClient(secretKey string) *StripeClient{
	stripe.Key = secretKey
	return &StripeClient{SecretKey: secretKey}
}

func (c *StripeClient) CreatePaymentIntent(amount int64, currency string) (*stripe.PaymentIntent, error){
	params := &stripe.PaymentIntentParams{
		Amount: stripe.Int64(amount),
		Currency: stripe.String(currency),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
			stripe.String("pix"),
		},
	}
	return paymentintent.New(params)
}
