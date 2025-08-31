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
	httpDialog "forum-dialog/internal/controller/http/handler/dialog"
)

func NewDialogClient(host, port string) *ClientDialog {
	return &ClientDialog{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/dialog",
	}
}

type ClientDialog struct {
	client  *http.Client
	baseURL string
}

func (c *ClientDialog) CreateDialog(
	account2ID int,
	accessToken string,
) (*httpDialog.CreateDialogResponse, error) {
	body := httpDialog.CreateDialogBody{
		Account2ID: account2ID,
	}

	response, err := c.postWithAccessToken("/create", body, accessToken)
	if err != nil {
		return nil, err
	}

	var createDialogResponse httpDialog.CreateDialogResponse
	err = json.Unmarshal(response, &createDialogResponse)
	if err != nil {
		return nil, err
	}

	return &createDialogResponse, err
}

func (c *ClientDialog) UploadFile(
	file []byte,
	fileName string,
	accessToken string,
) (*httpDialog.UploadFileResponse, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	formFile, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(formFile, bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	response, err := c.multipartFormDataWithAccessToken("/file/upload", writer.FormDataContentType(), accessToken, &body)
	if err != nil {
		return nil, err
	}

	var uploadFileResponse httpDialog.UploadFileResponse
	err = json.Unmarshal(response, &uploadFileResponse)
	if err != nil {
		return nil, err
	}

	return &uploadFileResponse, err
}

func (c *ClientDialog) DownloadFile(
	fileURL string,
) ([]byte, error) {
	response, err := c.get("/file/download/" + fileURL)
	if err != nil {
		return nil, err
	}

	return response, err
}

func (c *ClientDialog) MarkDialogAsStarred(
	dialogID int,
	accessToken string,
) error {
	_, err := c.postWithAccessToken("/mark/star/"+strconv.Itoa(dialogID), nil, accessToken)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientDialog) MarkMessagesAsRead(
	dialogID int,
	accessToken string,
) error {
	_, err := c.postWithAccessToken("/mark/read/"+strconv.Itoa(dialogID), nil, accessToken)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientDialog) DeleteDialog(
	dialogID int,
) error {
	_, err := c.post("/del/"+strconv.Itoa(dialogID), nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientDialog) DialogsByAccountID(
	accessToken string,
) (*httpDialog.DialogsByAccountIDResponse, error) {
	response, err := c.getWithAccessToken("/all", accessToken)
	if err != nil {
		return nil, err
	}

	var dialogsByAccountIDResponse httpDialog.DialogsByAccountIDResponse
	err = json.Unmarshal(response, &dialogsByAccountIDResponse)
	if err != nil {
		return nil, err
	}

	return &dialogsByAccountIDResponse, err
}

func (c *ClientDialog) MessagesByDialogID(
	dialogID int,
	accessToken string,
) (*httpDialog.MessagesByDialogIDResponse, error) {
	response, err := c.getWithAccessToken("/messages/"+strconv.Itoa(dialogID), accessToken)
	if err != nil {
		return nil, err
	}

	var messagesByDialogID httpDialog.MessagesByDialogIDResponse
	err = json.Unmarshal(response, &messagesByDialogID)
	if err != nil {
		return nil, err
	}

	return &messagesByDialogID, err
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

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ClientDialog) post(path string, body any) ([]byte, error) {
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

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ClientDialog) get(path string) ([]byte, error) {
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

func (c *ClientDialog) multipartFormDataWithAccessToken(path, contentType, accessToken string, body *bytes.Buffer) ([]byte, error) {

	req, err := http.NewRequest("POST", c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})

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
