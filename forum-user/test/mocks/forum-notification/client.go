package forum_notification

import (
	"context"
)

func New() *NotificationClient {
	return &NotificationClient{}
}

type NotificationClient struct{}

func (c *NotificationClient) NewWarningFromAdminNotification(ctx context.Context,
	toAccountID,
	adminAccountID int,
	warningText,
	adminLogin string,
) error {
	return nil
}
