package forum_user

import (
	"context"
)

func New() *UserClient {
	return &UserClient{}
}

type UserClient struct{}

func (c *UserClient) CreateUser(ctx context.Context,
	accountID int,
	login string,
) error {
	return nil
}
