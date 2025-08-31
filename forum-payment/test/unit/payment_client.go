package unit

import (
	"bytes"
	"encoding/json"
	"forum-payment/internal/model"
	"io"
	"net/http"
	"strconv"
)

import (
	httpPayment "forum-payment/internal/controller/http/handler/payment"
)

func NewPaymentClient(host, port string) *ClientPayment {
	return &ClientPayment{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/payment",
	}
}

type ClientPayment struct {
	client  *http.Client
	baseURL string
}

func (c *ClientPayment) CreatePayment(
	accountID int,
	productType model.ProductType,
	currency model.Currency,
	amountUSD float32,
) (*httpPayment.CreatePaymentResponse, error) {
	body := httpPayment.CreatePaymentBody{
		AccountID:   accountID,
		ProductType: productType,
		Currency:    currency,
		AmountUSD:   amountUSD,
	}

	response, err := c.post("/create", body)
	if err != nil {
		return nil, err
	}

	var createPaymentResponse httpPayment.CreatePaymentResponse
	err = json.Unmarshal(response, &createPaymentResponse)
	if err != nil {
		return nil, err
	}

	return &createPaymentResponse, err
}

func (c *ClientPayment) PaidPayment(
	paymentID int,
) error {
	_, err := c.post("/paid/"+strconv.Itoa(paymentID), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientPayment) CancelPayment(
	paymentID int,
) error {
	_, err := c.post("/cancel/"+strconv.Itoa(paymentID), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientPayment) StatusPayment(
	paymentID int,
) (*httpPayment.StatusPaymentResponse, error) {

	response, err := c.get("/status/" + strconv.Itoa(paymentID))
	if err != nil {
		return nil, err
	}

	var statusPaymentResponse httpPayment.StatusPaymentResponse
	err = json.Unmarshal(response, &statusPaymentResponse)
	if err != nil {
		return nil, err
	}

	return &statusPaymentResponse, err
}

func (c *ClientPayment) post(path string, body any) ([]byte, error) {
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

func (c *ClientPayment) get(path string) ([]byte, error) {
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
