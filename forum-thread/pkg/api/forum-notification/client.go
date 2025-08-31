package forum_notification

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func New(host, port string) *NotificationClient {
	return &NotificationClient{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/notification",
	}
}

type NotificationClient struct {
	client  *http.Client
	baseURL string
}

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

func (c *NotificationClient) post(path string, body map[string]any) ([]byte, error) {
	bytesBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(bytesBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}
