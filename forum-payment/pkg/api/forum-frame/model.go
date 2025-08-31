package forum_frame

type ConfirmPaymentForFrameBody struct {
	PaymentID            int    `json:"paymentID"`
	InterServerSecretKey string `json:"interServerSecretKey"`
}
