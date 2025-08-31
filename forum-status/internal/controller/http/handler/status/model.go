package status

import "forum-status/internal/model"

type CreatePaymentForStatusBody struct {
	StatusID int    `json:"statusID"`
	Duration int    `json:"duration"`
	Currency string `json:"currency"`
}

type CreatePaymentForStatusResponse struct {
	PaymentID int    `json:"paymentID"`
	Address   string `json:"address"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
}

type ConfirmPaymentForStatusBody struct {
	PaymentID            int    `json:"paymentID"`
	InterServerSecretKey string `json:"interServerSecretKey"`
}

type AssignStatusToAccountBody struct {
	AccountID            int    `json:"accountID"`
	StatusID             int    `json:"statusID"`
	InterServerSecretKey string `json:"interServerSecretKey"`
}

type RevokeStatusFromAccountBody struct {
	AccountID            int    `json:"accountID"`
	StatusID             int    `json:"statusID"`
	InterServerSecretKey string `json:"interServerSecretKey"`
}

type GetStatusByAccountIDResponse struct {
	Statuses []*model.Status `json:"statuses"`
}

type GetAllStatusResponse struct {
	Statuses []*model.StatusConst `json:"statuses"`
}
