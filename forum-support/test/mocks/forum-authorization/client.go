package forum_authorization

import (
	"context"
	"forum-support/internal/model"
	forum_authorization "forum-support/pkg/api/forum-authorization"

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
		TwoFaStatus: true,
		Message:     "ok",
		Code:        200,
	}

	if accessToken.Value == "user" {
		response.AccountID = 1
		response.Role = "user"
	} else if accessToken.Value == "support" {
		response.AccountID = 2
		response.Role = model.RoleSupport
	} else if accessToken.Value == "o_user" {
		response.AccountID = 3
		response.Role = "user"
	}

	return response, nil
}
