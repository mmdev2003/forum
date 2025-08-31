package message

import "forum-thread/internal/model"

type CreateTopicBody struct {
	SubthreadID     int    `json:"subthreadID" validate:"required"`
	ThreadID        int    `json:"threadID" validate:"required"`
	SubthreadName   string `json:"subthreadName" validate:"required"`
	ThreadName      string `json:"threadName" validate:"required"`
	TopicName       string `json:"topicName" validate:"required"`
	TopicOwnerLogin string `json:"topicOwnerLogin" validate:"required"`
}
type CreateTopicResponse struct {
	TopicID int `json:"topicID"`
}

type AddViewToTopicBody struct {
	TopicID int `json:"topicID" validate:"required"`
}

type CloseTopicBody struct {
	TopicOwnerAccountID int    `json:"topicOwnerAccountID" validate:"required"`
	AdminAccountID      int    `json:"adminAccountID" validate:"required"`
	TopicID             int    `json:"topicID" validate:"required"`
	TopicName           string `json:"topicName" validate:"required"`
	AdminLogin          string `json:"adminLogin" validate:"required"`
}

type ChangeTopicPriorityBody struct {
	SubthreadID   int `json:"subthreadID" validate:"required"`
	TopicID       int `json:"topicID" validate:"required"`
	TopicPriority int `json:"topicPriority" validate:"required"`
}

type UpdateTopicAvatarBody struct {
	TopicID int `json:"topicID" validate:"required"`
}

type TopicsBySubthreadIDResponse struct {
	Topics []*model.Topic `json:"topics"`
}

type TopicsByAccountIDResponse struct {
	Topics []*model.Topic `json:"topics"`
}

type TopicsOnModerationIDResponse struct {
	Topics []*model.Topic `json:"topics"`
}
