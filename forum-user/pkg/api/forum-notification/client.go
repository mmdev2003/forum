package forum_notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func New(host, port string) *NotificationClient {
	return &NotificationClient{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/notification",
	}
}

type NotificationClient struct {
	client  *http.Client
	baseURL string
}

func (c *NotificationClient) NewWarningFromAdminNotification(ctx context.Context,
	toAccountID,
	adminAccountID int,
	warningText,
	adminLogin string,
) error {
	return nil
}

func (c *NotificationClient) post(path string, body any) ([]byte, error) {
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
