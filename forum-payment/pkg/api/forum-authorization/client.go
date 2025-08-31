package forum_authorization

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func New(host, port string) *AuthorizationClient {
	return &AuthorizationClient{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/authorization",
	}
}

type AuthorizationClient struct {
	client  *http.Client
	baseURL string
}

func (c *AuthorizationClient) CheckAuthorization(ctx context.Context,
	request echo.Context,
) (*AuthorizationData, error) {
	accessToken, err := request.Cookie("Access-Token")
	if err != nil {
		return &AuthorizationData{
			AccountID:   0,
			Role:        "guest",
			TwoFaStatus: false,
			Message:     "guest",
			Code:        200,
		}, nil
	}

	response, err := c.getWithAccessToken("/check", accessToken.Value)
	if err != nil {
		return nil, err
	}

	var authorizationData AuthorizationData
	err = json.Unmarshal(response, &authorizationData)
	if err != nil {
		return nil, err
	}
	return &authorizationData, nil
}

func (c *AuthorizationClient) getWithAccessToken(path, accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, errors.New("status code not 200")
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}
