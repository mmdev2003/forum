package model

import "errors"

const (
	LoggerKey            = "logger"
	AuthorizationDataKey = "authorizationData"
)

const (
	RoleUser    = "user"
	RoleAdmin   = "admin"
	RoleSupport = "support"
)

const (
	ErrCheckAuthorizationFailed = "check authorization failed"
	ErrTokenExpired             = "token expired"
	ErrTokenInvalid             = "token invalid"
	ErrUnauthorized             = "Unauthorized"
)

var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrTwoFaAlreadyEnabled = errors.New("two-factor authentication already enabled")
	ErrTwoFaCodeInvalid    = errors.New("invalid two-factor authentication code")
	ErrTwoFaNotEnabled     = errors.New("two-factor authentication not enabled")
	ErrHttpRequestFailed   = errors.New("http request failed")
)

const (
	CodeErrAccessTokenExpired = 4012
	CodeErrAccessTokenInvalid = 4013
)

const PREFIX = "/api/authentication"
