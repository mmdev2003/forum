package forum_status

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

func New(host, port, interServerSecretKey string) *StatusClient {
	return &StatusClient{
		client:               &http.Client{},
		baseURL:              "http://" + host + ":" + port + "/api/status",
		interServerSecretKey: interServerSecretKey,
	}
}

type StatusClient struct {
	client               *http.Client
	baseURL              string
	interServerSecretKey string
}

func (c *StatusClient) StatusByAccountID(ctx context.Context,
	accountID int,
) ([]*Status, error) {
	response, err := c.get("/info/" + strconv.Itoa(accountID))
	if err != nil {
		return nil, err
	}

	var statusByAccountIDResponse StatusByAccountIDResponse
	err = json.Unmarshal(response, &statusByAccountIDResponse)
	if err != nil {
		return nil, err
	}

	return statusByAccountIDResponse.Statuses, nil
}

func (c *StatusClient) AssignStatusToAccount(ctx context.Context,
	statusID, accountID int,
) error {
	body := AssignStatusToAccountBody{
		AccountID:            accountID,
		StatusID:             statusID,
		InterServerSecretKey: c.interServerSecretKey,
	}
	_, err := c.post("/assign", body)
	if err != nil {
		return err
	}
	return nil
}

func (c *StatusClient) RevokeStatusFromAccount(ctx context.Context,
	statusID, accountID int,
) error {
	body := AssignStatusToAccountBody{
		AccountID:            accountID,
		StatusID:             statusID,
		InterServerSecretKey: c.interServerSecretKey,
	}
	_, err := c.post("/revoke", body)
	if err != nil {
		return err
	}
	return nil
}

func (c *StatusClient) post(path string, body any) ([]byte, error) {
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
		return nil, errors.New("status code is not 200")
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *StatusClient) get(path string) ([]byte, error) {
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
