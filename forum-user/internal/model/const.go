package model

var (
	BullingWarningType = "BULLING"
	CensureWarningType = "CENSURE"
)

var WarningMap = map[string]string{
	CensureWarningType: "Бан за цензуру",
	BullingWarningType: "Бан за буллинг",
}

var (
	RoleAdmin = "admin"
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

const (
	CodeErrAccessTokenExpired = 4012
	CodeErrAccessTokenInvalid = 4013
)

const (
	UserFullTextSearchIndex = "user"
)

const PREFIX = "/api/user"
