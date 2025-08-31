package forum_user

import (
	"context"
	"forum-thread/pkg/api/forum-user"
	"time"
)

func New() *UserClient {
	return &UserClient{}
}

type UserClient struct{}

func (c *UserClient) UserByLogin(ctx context.Context,
	login string,
) (*forum_user.GetUserByLoginResponse, error) {

	return &forum_user.GetUserByLoginResponse{
		User: []*forum_user.User{
			{
				ID:        1,
				AccountID: 1,

				Login:     login,
				HeaderUrl: "",
				AvatarUrl: "",

				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}, nil
}
