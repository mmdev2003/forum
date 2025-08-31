package forum_support

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func New(host, port string) *SupportClient {
	return &SupportClient{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/support",
	}
}

type SupportClient struct {
	client  *http.Client
	baseURL string
}

func (c *SupportClient) ReportMessage(ctx context.Context,
	accountID,
	messageID int,
	reportText string,
) error {
	body := ReportMessageBody{
		AccountID:  accountID,
		MessageID:  messageID,
		ReportText: reportText,
	}
	_, err := c.post("/message/report", body)
	if err != nil {
		return err
	}
	return nil
}

func (c *SupportClient) post(path string, body any) ([]byte, error) {
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
