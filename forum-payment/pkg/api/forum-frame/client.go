package forum_frame

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func New(host, port, interServerSecretKey string) *FrameClient {
	return &FrameClient{
		client:               &http.Client{},
		baseURL:              "http://" + host + ":" + port + "/api/frame",
		interServerSecretKey: interServerSecretKey,
	}
}

type FrameClient struct {
	client               *http.Client
	baseURL              string
	interServerSecretKey string
}

func (c *FrameClient) ConfirmPaymentForFrame(ctx context.Context,
	paymentID int,
) error {
	requestBody := ConfirmPaymentForFrameBody{
		PaymentID:            paymentID,
		InterServerSecretKey: c.interServerSecretKey,
	}

	_, err := c.post("/payment/confirm", requestBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *FrameClient) post(path string, body any) ([]byte, error) {
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
