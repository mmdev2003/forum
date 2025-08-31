package model

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

const PREFIX = "/api/frame"
