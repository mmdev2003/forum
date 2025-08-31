package forum_admin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func New(host, port, interServerSecretKey string) *AdminClient {
	return &AdminClient{
		client:               &http.Client{},
		baseURL:              "http://" + host + ":" + port + "/api/admin",
		interServerSecretKey: interServerSecretKey,
	}
}

type AdminClient struct {
	client               *http.Client
	baseURL              string
	interServerSecretKey string
}

func (c *AdminClient) CreateAdmin(ctx context.Context,
	accountID int,
) error {
	body := CreateAdminBody{
		AccountID:            accountID,
		InterServerSecretKey: c.interServerSecretKey,
	}

	_, err := c.post("/create", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *AdminClient) post(path string, body any) ([]byte, error) {
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
