package model

import "errors"

const (
	LoggerKey string = "logger"
)

var (
	ErrTokenExpired    = errors.New("token expired")
	ErrTokenInvalid    = errors.New("invalid token")
	ErrAccountNotFound = errors.New("account not found")
)

var (
	CodeErrAccessTokenNonSet  = 4011
	CodeErrAccessTokenExpired = 4012
	CodeErrAccessTokenInvalid = 4013
)

const PREFIX = "/api/authorization"
