package forum_authorization

import (
	"context"
	"forum-thread/internal/model"
	"forum-thread/pkg/api/forum-authorization"
	"github.com/labstack/echo/v4"
)

func New() *AuthorizationClient {
	return &AuthorizationClient{}
}

type AuthorizationClient struct{}

func (c *AuthorizationClient) CheckAuthorization(ctx context.Context,
	request echo.Context,
) (*forum_authorization.AuthorizationData, error) {
	accessToken, err := request.Cookie("Access-Token")
	if err != nil {
		return &forum_authorization.AuthorizationData{
			AccountID:   0,
			Role:        "guest",
			TwoFaStatus: false,
			Message:     "guest",
			Code:        200,
		}, nil
	}

	response := &forum_authorization.AuthorizationData{
		Role:        "user",
		TwoFaStatus: true,
		Message:     "ok",
		Code:        200,
	}

	if accessToken.Value == "user1" {
		response.AccountID = 1
	} else if accessToken.Value == "user2" {
		response.AccountID = 2
	} else if accessToken.Value == "user3" {
		response.AccountID = 3
	} else if accessToken.Value == "user4" {
		response.AccountID = 4
	} else if accessToken.Value == "admin1" {
		response.AccountID = 3
		response.Role = model.RoleAdmin
	}

	return response, nil
}
