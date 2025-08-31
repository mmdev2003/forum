package unit

import (
	"bytes"
	"encoding/json"
	"errors"
	support_request "forum-support/internal/controller/http/handler/support-request"
	"forum-support/internal/model"
	"io"
	"net/http"
	"strconv"
)

type SupportClient struct {
	client  *http.Client
	baseURL string
}

func (c *SupportClient) CreateSupportRequest(
	accessToken,
	title,
	description string,
) (int, error) {
	body := support_request.CreateSupportRequestRequest{
		Title:       title,
		Description: description,
	}
	rawResponse, err := c.postWithAuth("/request/create", accessToken, body)
	if err != nil {
		return -1, err
	}

	var response support_request.CreateSupportRequestResponse
	if err := json.Unmarshal(rawResponse, &response); err != nil {
		return -1, err
	}

	return response.SupportRequestID, nil
}

func (c *SupportClient) OpenSupportRequest(accessToken string, requestID int) error {
	_, err := c.postWithAuth("/request/open/"+strconv.Itoa(requestID), accessToken, nil)
	return err
}

func (c *SupportClient) CloseSupportRequest(accessToken string, requestID int) error {
	_, err := c.postWithAuth("/request/close/"+strconv.Itoa(requestID), accessToken, nil)
	return err
}

func (c *SupportClient) GetRequestById(accessToken string, requestID int) (*model.SupportRequest, error) {
	rawResponse, err := c.getWithAuth("/request/"+strconv.Itoa(requestID), accessToken)
	if err != nil {
		return nil, err
	}

	var response model.SupportRequest
	err = json.Unmarshal(rawResponse, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *SupportClient) GetRequests(accessToken string, requestID int) ([]model.SupportRequest, error) {
	rawResponse, err := c.getWithAuth("/request", accessToken)
	if err != nil {
		return nil, err
	}

	var response []model.SupportRequest
	err = json.Unmarshal(rawResponse, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func NewSupportClient(host, port string) *SupportClient {
	return &SupportClient{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/support",
	}
}

func (c *SupportClient) postWithAuth(path, accessToken string, body any) ([]byte, error) {
	bytesBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(bytesBody))
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

func (c *SupportClient) getWithAuth(path string, accessToken string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
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
