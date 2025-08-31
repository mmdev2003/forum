package forum_payment

type CreatePaymentBody struct {
	AccountID   int     `json:"accountID"`
	ProductType string  `json:"productType"`
	Currency    string  `json:"currency"`
	AmountUSD   float32 `json:"amountUSD"`
}

type PaymentData struct {
	PaymentID int    `json:"paymentID"`
	Amount    string `json:"amount"`
	Address   string `json:"address"`
	Currency  string `json:"currency"`
}
