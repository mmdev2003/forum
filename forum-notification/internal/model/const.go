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

const PREFIX string = "/api/notification"

type NotificationType string

const (
	MessageFromTopicType         NotificationType = "MessageFromTopic"
	MessageReplyFromTopicType    NotificationType = "MessageReplyFromTopic"
	LikeMessageFromTopicType     NotificationType = "LikeMessageFromTopic"
	TopicClosedType              NotificationType = "TopicClosed"
	ResponseToSupportRequestType NotificationType = "ResponseToSupportRequest"
	StatusReceivedType           NotificationType = "StatusReceived"
	FrameReceivedType            NotificationType = "FrameReceived"
	MessageFromDialogType        NotificationType = "MessageFromDialog"
	MentionFromTopicType         NotificationType = "MentionFromTopic"
	WarningFromAdminType         NotificationType = "WarningFromAdmin"
)
