package forum_payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func New(host, port string) *PaymentClient {
	return &PaymentClient{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/payment",
	}
}

type PaymentClient struct {
	client  *http.Client
	baseURL string
}

func (c *PaymentClient) CreatePayment(ctx context.Context,
	accountID int,
	productType string,
	currency string,
	amountUSD float32,
) (*PaymentData, error) {
	body := CreatePaymentBody{
		AccountID:   accountID,
		ProductType: productType,
		Currency:    currency,
		AmountUSD:   amountUSD,
	}

	response, err := c.post("/create", body)
	if err != nil {
		return nil, err
	}

	var paymentData PaymentData
	err = json.Unmarshal(response, &paymentData)
	if err != nil {
		return nil, err
	}
	return &paymentData, nil
}

func (c *PaymentClient) post(path string, body any) ([]byte, error) {
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
