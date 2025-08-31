package forum_payment

type PaymentData struct {
	PaymentID int    `json:"paymentID"`
	Amount    string `json:"amount"`
	Address   string `json:"address"`
	Currency  string `json:"currency"`
}
