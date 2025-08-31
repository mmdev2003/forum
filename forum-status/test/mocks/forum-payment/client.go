package forum_payment

import (
	"context"
	"forum-status/pkg/api/forum-payment"
	"net/http"
)

func New() *PaymentClient {
	return &PaymentClient{}
}

type PaymentClient struct {
	client  *http.Client
	baseURL string
}

func (c *PaymentClient) CreatePayment(ctx context.Context,
	accountID int,
	productType string,
	currency string,
	amountUSD float32,
) (*forum_payment.PaymentData, error) {
	return &forum_payment.PaymentData{
		PaymentID: 1,
		Amount:    "1",
		Address:   "fwewe",
		Currency:  "btc",
	}, nil
}
