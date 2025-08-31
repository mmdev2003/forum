package forum_notification

import (
	"forum-support/internal/model"
)

func New() *NotificationClient {
	return &NotificationClient{}
}

type NotificationClient struct{}

func (c NotificationClient) SendSupportRequestClosedNotification(requesterAccountID, supportRequestID int) error {
	return nil
}

func (c NotificationClient) SendStatusReceivedNotificationRequest(receiverAccountID int, statusName model.RequestStatus) error {
	return nil
}

func (c NotificationClient) SendResponseToSupportRequestNotificationRequest(requesterAccountID, supportRequestID int) error {
	return nil
}
