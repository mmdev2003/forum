package model

import "errors"

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
	ErrAccountStatisticNotFound = errors.New("account statistic not found")
)

const (
	CodeErrAccessTokenExpired = 4012
	CodeErrAccessTokenInvalid = 4013
)

const PREFIX = "/api/thread"

var (
	ErrNoPermissionToCreateTopicInThisThread    = errors.New("no permission to create topic in this thread")
	ErrNoPermissionToCreateTopicInThisSubthread = errors.New("no permission to create topic in this subthread")
	ErrMaxTopicsPerDay                          = errors.New("max topics per day")
	MaxTopicsInSubthread                        = errors.New("max topics in subthread")
)

const (
	AddViewToSubthreadQueue = "addViewToSubthread"

	CreateTopicQueue         = "createTopic"
	AddViewToTopicQueue      = "addViewToTopic"
	CloseTopicQueue          = "closeTopic"
	ChangeTopicPriorityQueue = "changeTopicPriority"

	SendMessageToTopicQueue = "sendMessageToTopic"
	LikeMessageQueue        = "likeMessage"
	ReportMessageQueue      = "reportMessage"
)

const (
	MessageFullTextSearchIndex = "message"
)

var (
	OnModerationTopicStatus = "onModeration"
	RejectedTopicStatus     = "rejected"
	ApprovedTopicStatus     = "approved"
)

var (
	RoleAdmin = "admin"
)

type StatusPermission struct {
	ID                   int
	PrivateThreads       []string
	PrivateSubthreads    []string
	PrivateTopics        []string
	MaxTopicsPerDay      int
	MaxTopicsInSubthread int
}

var StatusPermissionMap = map[int]StatusPermission{
	1: {
		1,
		[]string{"thread1", "thread2", "thread3", "thread4"},
		[]string{"subthread1", "subthread2", "subthread3", "subthread4"},
		[]string{"topic1", "topic2", "topic3", "topic4"},
		2,
		10,
	},
	2: {
		2,
		[]string{"thread1", "thread2", "thread3", "thread4"},
		[]string{"subthread1", "subthread2", "subthread3", "subthread4"},
		[]string{"topic1", "topic2", "topic3", "topic4"},
		2,
		10,
	},
	3: {
		3,
		[]string{"thread4"},
		[]string{"subthread4"},
		[]string{"topic4"},
		2,
		10,
	},
	4: {
		4,
		[]string{"thread1"},
		[]string{"subthread2"},
		[]string{"topic2"},
		2,
		10,
	},
}
