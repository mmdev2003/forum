package forum_status

import (
	"context"
)

func New() *StatusClient {
	return &StatusClient{}
}

type StatusClient struct{}

func (c *StatusClient) ConfirmPaymentForStatus(ctx context.Context,
	paymentID int,
) error {
	return nil
}
