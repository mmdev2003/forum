package model

type StatusConst struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	MonthlyPrice float32 `json:"monthlyPrice"`
	ForeverPrice float32 `json:"foreverPrice"`
}

var Statuses = map[int]StatusConst{
	1: {ID: 1, Name: "Agent", MonthlyPrice: 100, ForeverPrice: 200},
	2: {ID: 2, Name: "Reader", MonthlyPrice: 100, ForeverPrice: 200},
}

const (
	LoggerKey            = "logger"
	AuthorizationDataKey = "authorizationData"
)

type ProductType string

const (
	TStatus ProductType = "status"
)

type PaymentStatus string

const (
	Pending   PaymentStatus = "pending"
	Confirmed PaymentStatus = "confirmed"
	Canceled  PaymentStatus = "canceled"
)

const (
	ErrCheckAuthorizationFailed = "check authorization failed"
	ErrTokenExpired             = "token expired"
	ErrTokenInvalid             = "token invalid"
)

const (
	CodeErrAccessTokenExpired = 4012
	CodeErrAccessTokenInvalid = 4013
)

const PREFIX = "/api/status"
