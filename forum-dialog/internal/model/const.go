package model

const (
	LoggerKey            string = "logger"
	AuthorizationDataKey string = "authorizationData"
)

const (
	ErrCheckAuthorizationFailed        = "check authorization failed"
	ErrTokenExpired                    = "token expired"
	ErrTokenInvalid                    = "token invalid"
	ErrUnauthorized             string = "Unauthorized"
)

const (
	CodeErrAccessTokenExpired = 4012
	CodeErrAccessTokenInvalid = 4013
)

const PREFIX string = "/api/dialog"
