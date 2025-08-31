package forum_thread

import (
	"context"
)

func New() *ThreadClient {
	return &ThreadClient{}
}

type ThreadClient struct{}

func (c *ThreadClient) CreateAccountStatistic(ctx context.Context,
	accountID int,
) error {
	return nil
}
