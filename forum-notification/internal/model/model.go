package model

import "time"

type Notification struct {
	ID        int              `json:"id"`
	Type      NotificationType `json:"type"`
	IsRead    bool             `json:"isRead"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
}

type NewMessageFromTopicNotification struct {
	Notification
	SenderMessageID int `json:"senderMessageID"`
	SenderAccountID int `json:"senderAccountID"`
	TopicID         int `json:"topicID"`

	SenderMessageText string `json:"senderMessageText"`
	TopicName         string `json:"topicName"`
	SenderLogin       string `json:"senderLogin"`
}

type NewMessageReplyFromTopicNotification struct {
	Notification
	SenderMessageID int `json:"senderMessageID"`
	SenderAccountID int `json:"senderAccountID"`
	TopicID         int `json:"topicID"`

	SenderMessageText string `json:"senderMessageText"`
	TopicName         string `json:"topicName"`
	SenderLogin       string `json:"senderLogin"`
}

type NewLikeMessageFromTopicNotification struct {
	Notification
	SenderMessageID int `json:"senderMessageID"`
	LikerAccountID  int `json:"likerAccountID"`
	TopicID         int `json:"topicID"`

	SenderMessageText string `json:"senderMessageText"`
	TopicName         string `json:"topicName"`
	LikerLogin        string `json:"likerLogin"`
}

type TopicClosedNotification struct {
	Notification
	AdminAccountID int `json:"adminAccountID"`
	TopicID        int `json:"topicID"`

	TopicName  string `json:"topicName"`
	AdminLogin string `json:"adminLogin"`
}

type ResponseToSupportRequestNotification struct {
	Notification
	SupportRequestID int `json:"supportRequestID"`
}

type SupportRequestClosedNotification struct {
	Notification
	SupportRequestID int `json:"supportRequestID"`
}

type StatusReceivedNotification struct {
	Notification
	StatusName string `json:"statusName"`
}

type FrameReceivedNotification struct {
	Notification
	FrameName string `json:"frameName"`
}

type NewMessageFromDialogNotification struct {
	Notification
	MessageID       int `json:"message_id"`
	DialogID        int `json:"dialog_id"`
	SenderAccountID int `json:"sender_account_id"`

	MessageText string `json:"message_text"`
	SenderLogin string `json:"sender_login"`
}

type NewMentionFromTopicNotification struct {
	Notification
	MessageID       int `json:"senderMessageID"`
	SenderAccountID int `json:"senderAccountID"`

	SenderMessageText string `json:"senderMessageText"`
	TopicName         string `json:"topicName"`
	SenderLogin       string `json:"senderLogin"`
}

type NewWarningFromAdminNotification struct {
	Notification
	AdminAccountID int `json:"adminAccountID"`

	WarningText string `json:"warningText"`
	AdminLogin  string `json:"adminLogin"`
}
