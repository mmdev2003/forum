package unit

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

import (
	httpStatus "forum-status/internal/controller/http/handler/status"
)

func NewStatusClient(host, port string) *ClientStatus {
	return &ClientStatus{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/status",
	}
}

type ClientStatus struct {
	client  *http.Client
	baseURL string
}

func (c *ClientStatus) CreatePaymentForStatus(
	statusID,
	duration int,
	currency,
	accessToken string,
) (*httpStatus.CreatePaymentForStatusResponse, error) {
	body := httpStatus.CreatePaymentForStatusBody{
		StatusID: statusID,
		Duration: duration,
		Currency: currency,
	}

	response, err := c.postWithAccessToken("/payment/create", body, accessToken)
	if err != nil {
		return nil, err
	}

	var createPaymentForStatusResponse httpStatus.CreatePaymentForStatusResponse
	err = json.Unmarshal(response, &createPaymentForStatusResponse)
	if err != nil {
		return nil, err
	}

	return &createPaymentForStatusResponse, err
}

func (c *ClientStatus) ConfirmPaymentForStatus(
	paymentID int,
	interServerSecretKey string,
) error {
	body := httpStatus.ConfirmPaymentForStatusBody{
		PaymentID:            paymentID,
		InterServerSecretKey: interServerSecretKey,
	}

	_, err := c.post("/payment/confirm", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientStatus) AssignStatusToAccount(
	statusID,
	accountID int,
	interServerSecretKey string,
) error {
	body := httpStatus.AssignStatusToAccountBody{
		StatusID:             statusID,
		AccountID:            accountID,
		InterServerSecretKey: interServerSecretKey,
	}

	_, err := c.post("/assign", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientStatus) RevokeStatusFromAccount(
	statusID,
	accountID int,
	interServerSecretKey string,
) error {
	body := httpStatus.RevokeStatusFromAccountBody{
		StatusID:             statusID,
		AccountID:            accountID,
		InterServerSecretKey: interServerSecretKey,
	}

	_, err := c.post("/revoke", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientStatus) GetStatusByAccountID(
	accountID int,
) (*httpStatus.GetStatusByAccountIDResponse, error) {

	response, err := c.get("/info/" + strconv.Itoa(accountID))
	if err != nil {
		return nil, err
	}
	var getStatusByAccountIDResponse httpStatus.GetStatusByAccountIDResponse
	err = json.Unmarshal(response, &getStatusByAccountIDResponse)
	if err != nil {
		return nil, err
	}
	return &getStatusByAccountIDResponse, nil
}

func (c *ClientStatus) GetAllStatus() (*httpStatus.GetAllStatusResponse, error) {

	response, err := c.get("/all")
	if err != nil {
		return nil, err
	}
	var getAllStatusResponse httpStatus.GetAllStatusResponse
	err = json.Unmarshal(response, &getAllStatusResponse)
	if err != nil {
		return nil, err
	}
	return &getAllStatusResponse, nil
}

func (c *ClientStatus) postWithAccessToken(path string, body any, accessToken string) ([]byte, error) {
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

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ClientStatus) post(path string, body any) ([]byte, error) {
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

func (c *ClientStatus) getWithAccessToken(path, accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ClientStatus) get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

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
