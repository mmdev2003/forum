package forum_status

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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

func (c *StatusClient) ConfirmPaymentForStatus(ctx context.Context,
	paymentID int,
) error {
	requestBody := ConfirmPaymentForStatusBody{
		PaymentID:            paymentID,
		InterServerSecretKey: c.interServerSecretKey,
	}

	_, err := c.post("/payment/confirm", requestBody)
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

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(response))
	}

	return response, nil
}
