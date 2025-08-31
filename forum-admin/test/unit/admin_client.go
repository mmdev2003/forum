package unit

import (
	"bytes"
	"encoding/json"
	httpAdmin "forum-admin/internal/controller/http/handler/admin"
	"io"
	"net/http"
)

func NewAdminClient(host, port, interServerSecretKey string) *ClientAdmin {
	return &ClientAdmin{
		client:               &http.Client{},
		baseURL:              "http://" + host + ":" + port + "/api/admin",
		interServerSecretKey: interServerSecretKey,
	}
}

type ClientAdmin struct {
	client               *http.Client
	baseURL              string
	interServerSecretKey string
}

func (c *ClientAdmin) CreateAdmin(
	accountID int,
) error {
	body := httpAdmin.CreateAdminBody{
		AccountID:            accountID,
		InterServerSecretKey: c.interServerSecretKey,
	}

	_, err := c.post("/create", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientAdmin) AllAdmin() (*httpAdmin.AllAdminResponse, error) {
	response, err := c.get("/all")
	if err != nil {
		return nil, err
	}

	var allAdminResponse httpAdmin.AllAdminResponse
	err = json.Unmarshal(response, &allAdminResponse)
	if err != nil {
		return nil, err
	}

	return &allAdminResponse, nil
}

func (c *ClientAdmin) post(path string, body any) ([]byte, error) {
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

func (c *ClientAdmin) get(path string) ([]byte, error) {
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
