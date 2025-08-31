package unit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	handler "forum-notification/internal/api/http/handler/notification"
	"forum-notification/internal/model"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func NewNotificationClient(host, port string) *NotificationClient {
	return &NotificationClient{
		client:       &http.Client{},
		baseUrl:      "http://" + host + ":" + port + model.PREFIX,
		wsConnection: nil,
	}
}

type NotificationClient struct {
	client       *http.Client
	baseUrl      string
	wsConnection *websocket.Conn
}

func (c *NotificationClient) WebsocketConnect(accessToken string) error {
	err := c.WebsocketDisconnect()
	if err != nil {
		return err
	}

	headers := http.Header{}
	headers.Add("Cookie", fmt.Sprintf("Access-Token=%s", accessToken))

	wsUrl := "ws" + strings.TrimPrefix(c.baseUrl, "http") + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsUrl, headers)
	if err != nil {
		return err
	}

	c.wsConnection = ws
	return nil
}

func (c *NotificationClient) WebsocketDisconnect() error {
	if c.wsConnection != nil {
		c.wsConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

		err := c.wsConnection.Close()
		c.wsConnection = nil
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *NotificationClient) WebsocketReadMessage() ([]byte, error) {
	_, msg, err := c.wsConnection.ReadMessage()
	return msg, err
}

func (c *NotificationClient) GetNotifications(accessToken string) ([]json.RawMessage, error) {
	rawNotifications, err := c.getWithAuth("/", accessToken)
	if err != nil {
		return nil, err
	}

	var notifications []json.RawMessage
	err = json.Unmarshal(rawNotifications, &notifications)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (c *NotificationClient) getWithAuth(path string, accessToken string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, c.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}

	request.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(request)
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

func (c *NotificationClient) CreateNotification(notificationType model.NotificationType, notification any) (int, error) {
	jsonData, err := json.Marshal(notification)
	if err != nil {
		return 0, err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		c.baseUrl+"/create/"+string(notificationType),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return 0, err
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(request)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	var result handler.CreateNotificationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.NotificationID, nil
}

func (c *NotificationClient) GetFilters(accessToken string) (*handler.GetFiltersResponse, error) {
	rawFilters, err := c.getWithAuth("/filter", accessToken)
	if err != nil {
		return nil, err
	}

	var response handler.GetFiltersResponse
	err = json.Unmarshal(rawFilters, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *NotificationClient) SetFilters(accessToken string, enabledTypes []model.NotificationType) error {
	_, err := c.postWithAuth(
		"/filter",
		accessToken,
		handler.UpdateFiltersRequest{EnabledFilters: enabledTypes},
	)
	if err != nil {
		return nil
	}

	return nil
}

func (c *NotificationClient) postWithAuth(path, accessToken string, body any) ([]byte, error) {
	bytesBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", c.baseUrl+path, bytes.NewBuffer(bytesBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})

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
