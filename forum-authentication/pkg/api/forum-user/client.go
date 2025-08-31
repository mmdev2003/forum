package forum_user

import (
	"bytes"
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

func (c *UserClient) CreateUser(ctx context.Context,
	accountID int,
	login string,
) error {
	body := CreateUserBody{
		AccountID: accountID,
		Login:     login,
	}

	_, err := c.post("/create", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *UserClient) post(path string, body any) ([]byte, error) {
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
