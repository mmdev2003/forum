package unit

import (
	"bytes"
	"encoding/json"
	"forum-authorization/internal/model"
	"io"
	"net/http"
)

import (
	httpAuthorization "forum-authorization/internal/controller/http/handler/authorization"
)

func NewAuthorizationClient(host, port string) *ClientAuthorization {
	return &ClientAuthorization{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/authorization",
	}
}

type ClientAuthorization struct {
	client  *http.Client
	baseURL string
}

func (c *ClientAuthorization) Authorization(
	accountID int,
	role string,
	twoFaStatus bool,
) (*model.JWTTokens, error) {
	body := httpAuthorization.AuthBody{
		AccountID:   accountID,
		Role:        role,
		TwoFaStatus: twoFaStatus,
	}

	responseBody, err := c.post("/", body)
	if err != nil {
		return nil, err
	}

	var JWTTokens model.JWTTokens
	err = json.Unmarshal(responseBody, &JWTTokens)
	if err != nil {
		return nil, err
	}

	return &JWTTokens, err
}

func (c *ClientAuthorization) CheckAuthorization(
	accessToken string,
) (*httpAuthorization.CheckAuthorizationResponse, error) {
	responseBody, err := c.getWithAccessToken("/check", accessToken)
	if err != nil {
		return nil, err
	}

	var checkAuthorizationResponse httpAuthorization.CheckAuthorizationResponse
	err = json.Unmarshal(responseBody, &checkAuthorizationResponse)
	if err != nil {
		return nil, err
	}

	return &checkAuthorizationResponse, err
}

func (c *ClientAuthorization) RefreshTokens(
	refreshToken string,
) (*http.Response, error) {
	_, resp, err := c.getWithRefreshToken("/refresh", refreshToken)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *ClientAuthorization) post(path string, body any) ([]byte, error) {
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

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *ClientAuthorization) getWithAccessToken(path, accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return _body, nil
}

func (c *ClientAuthorization) getWithRefreshToken(path, refreshToken string) ([]byte, *http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, nil, err
	}
	req.AddCookie(&http.Cookie{Name: "Refresh-Token", Value: refreshToken})

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return _body, resp, nil
}
