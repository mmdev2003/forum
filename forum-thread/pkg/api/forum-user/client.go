package forum_user

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func New(host, port string) *UserClient {
	return &UserClient{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/user",
	}
}

type UserClient struct {
	client  *http.Client
	baseURL string
}

func (c *UserClient) UserByLogin(ctx context.Context,
	login string,
) (*GetUserByLoginResponse, error) {
	response, err := c.get("/info/login/" + login)
	if err != nil {
		return nil, err
	}

	var getUserByLoginResponse GetUserByLoginResponse
	err = json.Unmarshal(response, &getUserByLoginResponse)
	if err != nil {
		return nil, err
	}

	return &getUserByLoginResponse, nil
}

func (c *UserClient) get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
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
		return nil, errors.New("status code is not 200")
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}
