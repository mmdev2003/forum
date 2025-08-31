package unit

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	httpDialog "forum-support/internal/controller/http/handler/dialog"
	"forum-support/internal/model"
)

func NewDialogClient(host, port string) *ClientDialog {
	return &ClientDialog{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/support/dialog",
	}
}

type ClientDialog struct {
	client  *http.Client
	baseURL string
}

func (c *ClientDialog) CreateDialog(
	supportRequestID,
	userAccountID int,
	accessToken string,
) (int, error) {
	body := httpDialog.CreateDialogRequest{
		SupportRequestID: supportRequestID,
		UserAccountID:    userAccountID,
	}

	response, err := c.postWithAccessToken("/create", body, accessToken)
	if err != nil {
		return -1, err
	}

	var responseData struct {
		DialogID int `json:"dialogID"`
	}

	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return -1, err
	}

	return responseData.DialogID, err
}

func (c *ClientDialog) MarkMessagesAsRead(
	dialogID int,
	accessToken string,
) error {
	_, err := c.postWithAccessToken("/"+strconv.Itoa(dialogID)+"/mark/read", nil, accessToken)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientDialog) GetDialogs(
	accessToken string,
) ([]model.Dialog, error) {
	response, err := c.getWithAccessToken("/all", accessToken)
	if err != nil {
		return nil, err
	}

	var dialogs []model.Dialog
	err = json.Unmarshal(response, &dialogs)
	if err != nil {
		return nil, err
	}

	return dialogs, err
}

func (c *ClientDialog) MessagesByDialogID(
	dialogID int,
	accessToken string,
) ([]model.Message, error) {
	response, err := c.getWithAccessToken("/"+strconv.Itoa(dialogID), accessToken)
	if err != nil {
		return nil, err
	}

	var messages []model.Message
	err = json.Unmarshal(response, &messages)
	if err != nil {
		return nil, err
	}

	return messages, err
}

func (c *ClientDialog) postWithAccessToken(path string, body any, accessToken string) ([]byte, error) {
	bytesBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(bytesBody))
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})
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

func (c *ClientDialog) getWithAccessToken(path, accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		return nil, errors.New("status code is not 200")
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}
