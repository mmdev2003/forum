package model

import "errors"

type ProductType string

const (
	Status ProductType = "status"
	Frame  ProductType = "frame"
)

type PaymentStatus string

const (
	Pending   PaymentStatus = "pending"
	Confirmed PaymentStatus = "confirmed"
	Canceled  PaymentStatus = "canceled"
)

type Currency string

const (
	BTC Currency = "btc"
)

const (
	LoggerKey            = "logger"
	AuthorizationDataKey = "authorizationData"
)

const (
	ErrCheckAuthorizationFailed = "check authorization failed"
	ErrTokenExpired             = "token expired"
	ErrTokenInvalid             = "token invalid"
)

var (
	ErrPaymentNotFound = errors.New("payment not found")
)

const (
	CodeErrAccessTokenExpired = 4012
	CodeErrAccessTokenInvalid = 4013
)

const PREFIX = "/api/payment"
