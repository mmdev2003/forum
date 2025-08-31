package model

import "errors"

const (
	LoggerKey            string = "logger"
	AuthorizationDataKey string = "authorizationData"
)

const (
	ErrTitleTooLong       = "title too long"
	ErrDescriptionTooLong = "description too long"

	ErrCheckAuthorizationFailed        = "check authorization failed"
	ErrTokenExpired                    = "token expired"
	ErrTokenInvalid                    = "token invalid"
	ErrUnauthorized             string = "Unauthorized"
)

var ErrActionNotAllowed = errors.New("action is not allowed")

const (
	CodeErrAccessTokenExpired = 4012
	CodeErrAccessTokenInvalid = 4013
)

const PREFIX string = "/api/support"

type RequestStatus string

const (
	OpenRequestStatus   RequestStatus = "open"
	ClosedRequestStatus RequestStatus = "closed"
)

const (
	MaxTitleLength       = 150
	MaxDescriptionLength = 2500
)

const (
	RoleSupport = "support"
)
