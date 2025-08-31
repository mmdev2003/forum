package message

import "forum-thread/internal/model"

type SendMessageToTopicBody struct {
	SubthreadID                int    `json:"subthreadID" validate:"required"`
	TopicID                    int    `json:"topicID" validate:"required"`
	ReplyToMessageID           int    `json:"replyToMessageID"`
	ReplyMessageOwnerAccountID int    `json:"replyMessageOwnerAccountID"`
	TopicOwnerAccountID        int    `json:"topicOwnerAccountID" validate:"required"`
	SenderLogin                string `json:"senderLogin" validate:"required"`
	SenderMessageText          string `json:"senderMessageText" validate:"required"`
	ThreadName                 string `json:"threadName" validate:"required"`
	SubthreadName              string `json:"subthreadName" validate:"required"`
	TopicName                  string `json:"topicName" validate:"required"`
}

type SendMessageToTopicResponse struct {
	FilesURLs []string `json:"filesURLs"`
}

type LikeMessageBody struct {
	TopicID               int    `json:"topicID" validate:"required"`
	MessageOwnerAccountID int    `json:"messageOwnerAccountID" validate:"required"`
	LikerAccountID        int    `json:"likerAccountID" validate:"required"`
	LikeMessageID         int    `json:"likeMessageID" validate:"required"`
	LikeTypeID            int    `json:"likeTypeID" validate:"required"`
	LikerLogin            string `json:"likerLogin" validate:"required"`
	TopicName             string `json:"topicName" validate:"required"`
	LikeMessageText       string `json:"likeMessageText" validate:"required"`
}

type UnlikeMessageBody struct {
	LikeMessageID int `json:"likeMessageID" validate:"required"`
}
type ReportMessageBody struct {
	MessageID  int    `json:"messageID" validate:"required"`
	AccountID  int    `json:"accountID" validate:"required"`
	ReportText string `json:"reportText" validate:"required"`
}

type EditMessageBody struct {
	MessageID   int    `json:"messageID" validate:"required"`
	MessageText string `json:"messageText" validate:"required"`
}

type MessagesByTopicIDResponse struct {
	Messages []*model.Message `json:"messages"`
	Likes    []*model.Like    `json:"likes"`
	Files    []*model.File    `json:"files"`
}

type MessagesByAccountIDResponse struct {
	Messages []*model.Message `json:"messages"`
	Files    []*model.File    `json:"files"`
}

type MessagesByTextResponse struct {
	Messages []*model.MessageSearch `json:"messages"`
}
