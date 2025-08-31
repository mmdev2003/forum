package forum_support

import (
	"context"
)

func New() *MockSupportClient {
	return &MockSupportClient{}
}

type MockSupportClient struct{}

func (c *MockSupportClient) ReportMessage(ctx context.Context,
	accountID,
	messageID int,
	reportText string,
) error {
	return nil
}
