package forum_admin

import (
	"context"
)

func New() *AdminClient {
	return &AdminClient{}
}

type AdminClient struct{}

func (c *AdminClient) CreateAdmin(ctx context.Context,
	accountID int,
) error {

	return nil
}
