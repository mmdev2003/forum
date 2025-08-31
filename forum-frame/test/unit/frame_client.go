package unit

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

import (
	httpFrame "forum-frame/internal/controller/http/handler/frame"
)

func NewFrameClient(host, port string) *ClientFrame {
	return &ClientFrame{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/frame",
	}
}

type ClientFrame struct {
	client  *http.Client
	baseURL string
}

func (c *ClientFrame) CreatePaymentForFrame(
	frameID,
	duration int,
	currency,
	accessToken string,
) (*httpFrame.CreatePaymentForFrameResponse, error) {
	body := httpFrame.CreatePaymentForFrameBody{
		FrameID:  frameID,
		Duration: duration,
		Currency: currency,
	}

	response, err := c.postWithAccessToken("/payment/create", body, accessToken)
	if err != nil {
		return nil, err
	}

	var createPaymentForFrameResponse httpFrame.CreatePaymentForFrameResponse
	err = json.Unmarshal(response, &createPaymentForFrameResponse)
	if err != nil {
		return nil, err
	}

	return &createPaymentForFrameResponse, err
}

func (c *ClientFrame) AddNewFrame(
	frameFile []byte,
	monthlyPrice,
	foreverPrice float64,
	name string,
) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	formFile, err := writer.CreateFormFile("frameFile", "example_frame.png")
	if err != nil {
		return err
	}
	_, err = io.Copy(formFile, bytes.NewReader(frameFile))
	if err != nil {
		return err
	}

	err = writer.WriteField("monthlyPrice", strconv.FormatFloat(monthlyPrice, 'f', -1, 64))
	if err != nil {
		return err
	}

	err = writer.WriteField("foreverPrice", strconv.FormatFloat(foreverPrice, 'f', -1, 64))
	if err != nil {
		return err
	}

	err = writer.WriteField("name", name)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	_, err = c.multipartFormData("/new", writer.FormDataContentType(), &requestBody)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientFrame) ConfirmPaymentForFrame(
	paymentID int,
	interServerSecretKey string,
) error {
	body := httpFrame.ConfirmPaymentForFrameBody{
		PaymentID:            paymentID,
		InterServerSecretKey: interServerSecretKey,
	}

	_, err := c.post("/payment/confirm", "application/json", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientFrame) ChangeCurrentFrame(
	dbFrameID int,
	accessToken string,
) error {

	_, err := c.postWithAccessToken("/change/"+strconv.Itoa(dbFrameID), nil, accessToken)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientFrame) FramesByAccountID(
	accountID int,
) (*httpFrame.FramesByAccountIDResponse, error) {

	response, err := c.get("/" + strconv.Itoa(accountID))
	if err != nil {
		return nil, err
	}
	var framesByAccountIDResponse httpFrame.FramesByAccountIDResponse
	err = json.Unmarshal(response, &framesByAccountIDResponse)
	if err != nil {
		return nil, err
	}
	return &framesByAccountIDResponse, nil
}

func (c *ClientFrame) AllFrame() (*httpFrame.AllFrameResponse, error) {

	response, err := c.get("/all")
	if err != nil {
		return nil, err
	}
	var allFrameResponse httpFrame.AllFrameResponse
	err = json.Unmarshal(response, &allFrameResponse)
	if err != nil {
		return nil, err
	}
	return &allFrameResponse, nil
}

func (c *ClientFrame) DownloadFrame(
	frameID int,
) ([]byte, error) {
	frameFile, err := c.get("/download/" + strconv.Itoa(frameID))
	if err != nil {
		return nil, err
	}
	return frameFile, err
}

func (c *ClientFrame) postWithAccessToken(path string, body any, accessToken string) ([]byte, error) {
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

func (c *ClientFrame) post(path, contentType string, body any) ([]byte, error) {
	bytesBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(bytesBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

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

func (c *ClientFrame) multipartFormData(path string, contentType string, body *bytes.Buffer) ([]byte, error) {

	req, err := http.NewRequest("POST", c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

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

func (c *ClientFrame) getWithAccessToken(path, accessToken string) ([]byte, error) {
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

func (c *ClientFrame) get(path string) ([]byte, error) {
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
