package payment

import "forum-payment/internal/model"

type CreatePaymentBody struct {
	AccountID   int               `json:"accountID"`
	ProductType model.ProductType `json:"productType"`
	Currency    model.Currency    `json:"currency"`
	AmountUSD   float32           `json:"amountUSD"`
}

type CreatePaymentResponse struct {
	PaymentID int            `json:"paymentID"`
	Amount    string         `json:"amount"`
	Address   string         `json:"address"`
	Currency  model.Currency `json:"currency"`
}

type StatusPaymentResponse struct {
	PaymentStatus model.PaymentStatus `json:"paymentStatus"`
}
