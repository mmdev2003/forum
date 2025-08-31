package forum_authorization

import (
	"bytes"
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

func (c *AuthorizationClient) Authorization(ctx context.Context,
	accountID int,
) (*JWTTokens, error) {
	requestBody := map[string]any{
		"accountID": accountID,
	}

	response, err := c.post("/", requestBody)
	if err != nil {
		return nil, err
	}

	var jwtTokens JWTTokens
	err = json.Unmarshal(response, &jwtTokens)
	if err != nil {
		return nil, err
	}
	return &jwtTokens, nil
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

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *AuthorizationClient) post(path string, body map[string]any) ([]byte, error) {
	bytesBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(bytesBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, errors.New("status code not 200")
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}
