package notification

import (
	"encoding/json"
	"forum-notification/internal/model"
	"time"
)

type NotificationWsBody struct {
	ID        int                    `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	Type      model.NotificationType `json:"type"`
	Payload   json.RawMessage        `json:"payload"`
}

type MessageFromTopicNotificationPayload struct {
	MessageID        int `json:"message_id"`
	ReplierAccountID int `json:"replier_account_id"`
	TopicID          int `json:"topic_id"`

	MessageText  string `json:"message_text"`
	TopicName    string `json:"topic_name"`
	ReplierLogin string `json:"replied_login"`
}

type NewMessageReplyFromTopicNotificationPayload struct {
	MessageID        int `json:"message_id"`
	ReplierAccountID int `json:"replier_account_id"`
	TopicID          int `json:"topic_id"`

	MessageText  string `json:"message_text"`
	TopicName    string `json:"topic_name"`
	ReplierLogin string `json:"replied_login"`
}

type NewLikeMessageFromTopicNotificationPayload struct {
	MessageID      int `json:"message_id"`
	LikerAccountID int `json:"liker_account_id"`
	TopicID        int `json:"topic_id"`

	MessageText string `json:"message_text"`
	TopicName   string `json:"topic_name"`
	LikerLogin  string `json:"liker_login"`
}

type TopicClosedNotificationPayload struct {
	AdminAccountID int `json:"admin_account_id"`
	TopicID        int `json:"topic_id"`

	TopicName  string `json:"topic_name"`
	AdminLogin string `json:"admin_login"`
}

type ResponseToSupportRequestNotificationPayload struct {
	SupportRequestID int `json:"support_request_id"`
}

type SupportRequestClosedNotificationPayload struct {
	SupportRequestID int `json:"support_request_id"`
}

type StatusReceivedNotificationPayload struct {
	StatusName string `json:"status_name"`
}

type FrameReceivedNotificationPayload struct {
	FrameName string `json:"frame_name"`
}

type NewMessageFromDialogNotificationPayload struct {
	MessageID       int `json:"message_id"`
	DialogID        int `json:"dialog_id"`
	SenderAccountID int `json:"sender_account_id"`

	MessageText string `json:"message_text"`
	SenderLogin string `json:"sender_login"`
}

type NewMentionFromTopicNotificationPayload struct {
	MessageID        int `json:"message_id"`
	MentionAccountID int `json:"mention_account_id"`

	MessageText  string `json:"message_text"`
	TopicName    string `json:"topic_name"`
	MentionLogin string `json:"mention_login"`
}

type NewWarningFromAdminNotificationPayload struct {
	AdminAccountID int `json:"admin_account_id"`

	WarningText string `json:"message_text"`
	AdminLogin  string `json:"admin_login"`
}
