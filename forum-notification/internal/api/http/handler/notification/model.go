package notification

import "forum-notification/internal/model"

type GetFiltersResponse struct {
	EnabledFilters []model.NotificationType `json:"enabledFilters"`
}

type UpdateFiltersRequest struct {
	EnabledFilters []model.NotificationType `json:"enabledFilters"`
}

type CreateNotificationResponse struct {
	NotificationID int `json:"notificationID"`
}

type (
	NewMessageFromTopicNotificationRequest struct {
		TopicOwnerAccountID int `json:"topicOwnerAccountID"`
		SenderMessageID     int `json:"senderMessageID"`
		SenderAccountID     int `json:"senderAccountID"`
		TopicID             int `json:"topicID"`

		SenderMessageText string `json:"senderMessageText"`
		TopicName         string `json:"topicName"`
		SenderLogin       string `json:"senderLogin"`
	}

	NewMessageReplyFromTopicNotificationRequest struct {
		ReplyMessageOwnerAccountID int `json:"replyMessageOwnerAccountID"`
		SenderMessageID            int `json:"senderMessageID"`
		SenderAccountID            int `json:"senderAccountID"`
		TopicID                    int `json:"topicID"`

		SenderMessageText string `json:"senderMessageText"`
		TopicName         string `json:"topicName"`
		SenderLogin       string `json:"senderLogin"`
	}

	NewLikeMessageFromTopicNotificationRequest struct {
		MessageOwnerAccountID int `json:"messageOwnerAccountID"`
		SenderMessageID       int `json:"senderMessageID"`
		LikerAccountID        int `json:"likerAccountID"`
		TopicID               int `json:"topicID"`

		SenderMessageText string `json:"senderMessageText"`
		TopicName         string `json:"topicName"`
		LikerLogin        string `json:"likerLogin"`
	}

	TopicClosedNotificationRequest struct {
		TopicOwnerAccountID int `json:"topicOwnerAccountID"`
		AdminAccountID      int `json:"adminAccountID"`
		TopicID             int `json:"topicID"`

		TopicName  string `json:"topicName"`
		AdminLogin string `json:"adminLogin"`
	}

	ResponseToSupportRequestNotificationRequest struct {
		RequesterAccountID int `json:"requesterAccountID"`
		SupportRequestID   int `json:"supportRequestID"`
	}

	SupportRequestClosedNotificationRequest struct {
		RequesterAccountID int `json:"requesterAccountID"`
		SupportRequestID   int `json:"supportRequestID"`
	}

	StatusReceivedNotificationRequest struct {
		ReceiverAccountID int    `json:"receiverAccountID"`
		StatusName        string `json:"statusMame"`
	}

	FrameReceivedNotificationRequest struct {
		ReceiverAccountID int    `json:"receiverAccountID"`
		FrameName         string `json:"frameName"`
	}

	NewMessageFromDialogNotificationRequest struct {
		AccountID       int `json:"account_id"`
		MessageID       int `json:"message_id"`
		DialogID        int `json:"dialog_id"`
		SenderAccountID int `json:"sender_account_id"`

		MessageText string `json:"message_text"`
		SenderLogin string `json:"sender_login"`
	}

	NewMentionFromTopicNotificationRequest struct {
		MentionedAccountID int `json:"mentionedAccountID"`
		MessageID          int `json:"senderMessageID"`
		SenderAccountID    int `json:"senderAccountID"`

		SenderMessageText string `json:"senderMessageText"`
		TopicName         string `json:"topicName"`
		SenderLogin       string `json:"senderLogin"`
	}

	NewWarningFromAdminNotificationRequest struct {
		AccountID      int `json:"warnedAccountID"`
		AdminAccountID int `json:"adminAccountID"`

		WarningText string `json:"warningText"`
		AdminLogin  string `json:"adminLogin"`
	}
)
