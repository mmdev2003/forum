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
	httpUser "forum-user/internal/controller/http/handler/user"
)

func NewUserClient(host, port string) *ClientUser {
	return &ClientUser{
		client:  &http.Client{},
		baseURL: "http://" + host + ":" + port + "/api/user",
	}
}

type ClientUser struct {
	client  *http.Client
	baseURL string
}

func (c *ClientUser) BanUser(
	toAccountID int,
	accessToken string,
) error {
	body := httpUser.BanUserBody{
		ToAccountID: toAccountID,
	}

	_, err := c.postWithAccessToken("/ban/create", body, accessToken)
	if err != nil {
		return nil
	}

	return nil
}

func (c *ClientUser) UnbanUser(
	toAccountID int,
	accessToken string,
) error {

	_, err := c.postWithAccessToken("/ban/del/"+strconv.Itoa(toAccountID), nil, accessToken)
	if err != nil {
		return nil
	}

	return nil
}

func (c *ClientUser) CreateUser(
	accountID int,
	login string,
) error {
	body := httpUser.CreateUserBody{
		AccountID: accountID,
		Login:     login,
	}

	_, err := c.post("/create", body)
	if err != nil {
		return nil
	}

	return nil
}

func (c *ClientUser) NewWarningFromAdmin(
	toAccountID int,
	warningType,
	accessToken,
	adminLogin string,
) error {
	body := httpUser.NewWarningFromAdminBody{
		ToAccountID: toAccountID,
		WarningType: warningType,
		AdminLogin:  adminLogin,
	}

	_, err := c.postWithAccessToken("/warning/create", body, accessToken)
	if err != nil {
		return nil
	}

	return nil
}

func (c *ClientUser) UploadAvatar(
	avatarFile []byte,
	accountID int,
	accessToken string,
) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	formFile, err := writer.CreateFormFile("avatar", "example_frame.png")
	if err != nil {
		return err
	}
	_, err = io.Copy(formFile, bytes.NewReader(avatarFile))
	if err != nil {
		return err
	}

	err = writer.WriteField("accountID", strconv.Itoa(accountID))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	_, err = c.multipartFormDataWithAccessToken("/avatar/upload", writer.FormDataContentType(), accessToken, &requestBody)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientUser) UploadHeader(
	avatarFile []byte,
	accountID int,
	accessToken string,
) error {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	formFile, err := writer.CreateFormFile("header", "example_frame.png")
	if err != nil {
		return err
	}
	_, err = io.Copy(formFile, bytes.NewReader(avatarFile))
	if err != nil {
		return err
	}

	err = writer.WriteField("accountID", strconv.Itoa(accountID))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	_, err = c.multipartFormDataWithAccessToken("/header/upload", writer.FormDataContentType(), accessToken, &requestBody)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientUser) DownloadAvatar(
	fileID string,
) ([]byte, error) {
	avatarFile, err := c.get("/avatar/download/" + fileID)
	if err != nil {
		return nil, err
	}
	return avatarFile, err
}

func (c *ClientUser) DownloadHeader(
	fileID string,
) ([]byte, error) {
	avatarFile, err := c.get("/header/download/" + fileID)
	if err != nil {
		return nil, err
	}
	return avatarFile, err
}

func (c *ClientUser) GetUserByAccountID(
	accountID int,
) (*httpUser.GetUserByAccountIDResponse, error) {
	response, err := c.get("/info/id/" + strconv.Itoa(accountID))
	if err != nil {
		return nil, err
	}

	var getUserByAccountIDResponse httpUser.GetUserByAccountIDResponse
	err = json.Unmarshal(response, &getUserByAccountIDResponse)
	if err != nil {
		return nil, err
	}
	return &getUserByAccountIDResponse, err
}

func (c *ClientUser) UsersByLoginSearch(
	login string,
) (*httpUser.UsersByLoginResponse, error) {
	response, err := c.get("/search/login/" + login)
	if err != nil {
		return nil, err
	}

	var usersByLoginResponse httpUser.UsersByLoginResponse
	err = json.Unmarshal(response, &usersByLoginResponse)
	if err != nil {
		return nil, err
	}
	return &usersByLoginResponse, err
}

func (c *ClientUser) UserByLogin(
	login string,
) (*httpUser.GetUserByLoginResponse, error) {
	response, err := c.get("/info/login/" + login)
	if err != nil {
		return nil, err
	}

	var getUserByLoginResponse httpUser.GetUserByLoginResponse
	err = json.Unmarshal(response, &getUserByLoginResponse)
	if err != nil {
		return nil, err
	}
	return &getUserByLoginResponse, err
}

func (c *ClientUser) BanByAccountID(
	accessToken string,
) (*httpUser.BanByAccountIDResponse, error) {
	response, err := c.getWithAccessToken("/info/bans", accessToken)
	if err != nil {
		return nil, err
	}

	var banByAccountIDResponse httpUser.BanByAccountIDResponse
	err = json.Unmarshal(response, &banByAccountIDResponse)
	if err != nil {
		return nil, err
	}
	return &banByAccountIDResponse, err
}

func (c *ClientUser) AllWarningFromAdmin(
	toAccountID int,
) (*httpUser.AllWarningFromAdminResponse, error) {
	response, err := c.get("/warning/all/" + strconv.Itoa(toAccountID))
	if err != nil {
		return nil, err
	}

	var allWarningFromAdminResponse httpUser.AllWarningFromAdminResponse
	err = json.Unmarshal(response, &allWarningFromAdminResponse)
	if err != nil {
		return nil, err
	}
	return &allWarningFromAdminResponse, err
}

func (c *ClientUser) multipartFormData(path string, contentType string, body *bytes.Buffer) ([]byte, error) {

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

func (c *ClientUser) multipartFormDataWithAccessToken(path, contentType, accessToken string, body *bytes.Buffer) ([]byte, error) {

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

func (c *ClientUser) postWithAccessToken(path string, body any, accessToken string) ([]byte, error) {
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

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ClientUser) post(path string, body any) ([]byte, error) {
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

func (c *ClientUser) getWithAccessToken(path, accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
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

func (c *ClientUser) get(path string) ([]byte, error) {
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
