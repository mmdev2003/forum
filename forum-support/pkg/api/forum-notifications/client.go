package forum_notifications

import (
	"bytes"
	"encoding/json"
	"forum-support/internal/model"
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

func (c *NotificationClient) SendResponseToSupportRequestNotificationRequest(requesterAccountID, supportRequestID int) error {
	_, err := c.post("/create/StatusReceived", ResponseToSupportRequestNotificationRequest{
		RequesterAccountID: requesterAccountID,
		SupportRequestID:   supportRequestID,
	})
	return err
}

func (c *NotificationClient) SendStatusReceivedNotificationRequest(receiverAccountID int, statusName model.RequestStatus) error {
	_, err := c.post("/create/StatusReceived", StatusReceivedNotificationRequest{
		ReceiverAccountID: receiverAccountID,
		StatusName:        string(statusName),
	})
	return err
}

func (c *NotificationClient) SendSupportRequestClosedNotification(requesterAccountID, supportRequestID int) error {
	_, err := c.post("/create/TopicClosed", SupportRequestClosedNotificationRequest{
		RequesterAccountID: requesterAccountID,
		SupportRequestID:   supportRequestID,
	})
	return err
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

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}
