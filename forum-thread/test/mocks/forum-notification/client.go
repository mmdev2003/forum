package forum_notification

import (
	"context"
)

func New() *NotificationClient {
	return &NotificationClient{}
}

type NotificationClient struct{}

func (c *NotificationClient) NewMessageFromTopicNotification(ctx context.Context,
	topicOwnerAccountID,
	senderMessageID,
	senderAccountID,
	topicID int,
	senderMessageText,
	TopicName,
	senderLogin string,
) error {
	return nil
}

func (c *NotificationClient) NewMentionFromTopicNotification(ctx context.Context,
	mentionedAccountID,
	senderMessageID,
	senderAccountID int,
	senderMessageText,
	topicName,
	senderLogin string,
) error {
	return nil
}

func (c *NotificationClient) NewLikeMessageFromTopicNotification(ctx context.Context,
	messageOwnerAccountID,
	senderMessageID,
	likerAccountID,
	topicID int,
	senderMessageText,
	topicName,
	likerLogin string,
) error {
	return nil
}

func (c *NotificationClient) NewMessageReplyFromTopicNotification(ctx context.Context,
	replyMessageOwnerAccountID,
	senderMessageID,
	senderAccountID,
	topicID int,
	senderMessageText,
	topicName,
	senderLogin string,
) error {
	return nil
}

func (c *NotificationClient) TopicClosedNotificationRequest(ctx context.Context,
	topicOwnerAccountID,
	adminAccountID,
	topicID int,
	topicName,
	adminLogin string,
) error {
	return nil
}
