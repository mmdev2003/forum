package frame

import "forum-frame/internal/model"

type CreatePaymentForFrameBody struct {
	FrameID  int    `json:"frameID" validate:"required"`
	Duration int    `json:"duration" validate:"required"`
	Currency string `json:"currency" validate:"required"`
}

type CreatePaymentForFrameResponse struct {
	PaymentID int    `json:"paymentID"`
	Address   string `json:"address"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
}

type ConfirmPaymentForFrameBody struct {
	PaymentID            int    `json:"paymentID" validate:"required"`
	InterServerSecretKey string `json:"interServerSecretKey" validate:"required"`
}

type FramesByAccountIDResponse struct {
	Frames       []*model.Frame      `json:"frames"`
	CurrentFrame *model.CurrentFrame `json:"currentFrame"`
}

type AllFrameResponse struct {
	Frames []*model.FrameData `json:"frames"`
}
