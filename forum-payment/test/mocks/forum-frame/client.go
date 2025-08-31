package forum_frame

import (
	"context"
)

func New() *FrameClient {
	return &FrameClient{}
}

type FrameClient struct{}

func (c *FrameClient) ConfirmPaymentForFrame(ctx context.Context,
	paymentID int,
) error {
	return nil
}
