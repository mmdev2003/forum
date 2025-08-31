package forum_status

type ConfirmPaymentForStatusBody struct {
	PaymentID            int    `json:"paymentID"`
	InterServerSecretKey string `json:"interServerSecretKey"`
}
