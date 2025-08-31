package unit

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

import (
	httpAuthentication "forum-authentication/internal/controller/http/handler/authentication"
)

func NewAuthenticationClient(host, port, interServerSecretKey string) *ClientAuthentication {
	return &ClientAuthentication{
		client:               &http.Client{},
		baseURL:              "http://" + host + ":" + port + "/api/authentication",
		interServerSecretKey: interServerSecretKey,
	}
}

type ClientAuthentication struct {
	client               *http.Client
	baseURL              string
	interServerSecretKey string
}

func (c *ClientAuthentication) Register(
	login,
	email,
	password string,
) (*httpAuthentication.RegisterResponse, *http.Response, error) {
	body := httpAuthentication.RegisterBody{
		Login:    login,
		Email:    email,
		Password: password,
	}

	responseBody, resp, err := c.post("/register", body)
	if err != nil {
		return nil, nil, err
	}

	var registerResponse httpAuthentication.RegisterResponse
	err = json.Unmarshal(responseBody, &registerResponse)
	if err != nil {
		return nil, nil, err
	}

	return &registerResponse, resp, err
}

func (c *ClientAuthentication) Login(
	login,
	password,
	twoFaCode string,
) (*httpAuthentication.LoginResponse, *http.Response, error) {
	body := httpAuthentication.LoginBody{
		Login:     login,
		Password:  password,
		TwoFaCode: twoFaCode,
	}

	responseBody, resp, err := c.post("/login", body)
	if err != nil {
		return nil, nil, err
	}

	var loginResponse httpAuthentication.LoginResponse
	err = json.Unmarshal(responseBody, &loginResponse)
	if err != nil {
		return nil, nil, err
	}

	return &loginResponse, resp, err
}

func (c *ClientAuthentication) GenerateTwoFa(
	accessToken string,
) ([]byte, string, error) {

	req, err := http.NewRequest("GET", c.baseURL+"/2fa/generate", nil)
	if err != nil {
		return nil, "", err
	}
	req.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	twoFaKey := resp.Header.Get("Two-Fa-Key")

	qrImage, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return qrImage, twoFaKey, nil
}

func (c *ClientAuthentication) SetTwoFa(
	twoFaKey,
	twoFaCode,
	accessToken string,
) error {
	body := httpAuthentication.SetTwoFaBody{
		TwoFaKey:  twoFaKey,
		TwoFaCode: twoFaCode,
	}

	_, _, err := c.postWithAccessToken("/2fa/set", accessToken, body)
	if err != nil {
		return nil
	}

	return nil
}

func (c *ClientAuthentication) DeleteTwoFaKey(
	twoFaCode,
	accessToken string,
) error {

	_, _, err := c.postWithAccessToken("/2fa/del/"+twoFaCode, accessToken, nil)
	if err != nil {
		return nil
	}

	return nil
}

func (c *ClientAuthentication) VerifyTwoFa(
	twoFaCode,
	accessToken string,
) (*httpAuthentication.VerifyTwoFaResponse, error) {

	response, err := c.getWithAccessToken("/2fa/verify/"+twoFaCode, accessToken)
	if err != nil {
		return nil, err
	}

	var verifyTwoFaResponse httpAuthentication.VerifyTwoFaResponse
	err = json.Unmarshal(response, &verifyTwoFaResponse)
	if err != nil {
		return nil, err
	}

	return &verifyTwoFaResponse, err
}

func (c *ClientAuthentication) UpgradeToAdmin(
	accountID int,
) error {
	body := httpAuthentication.UpgradeToAdminBody{
		AccountID:            accountID,
		InterServerSecretKey: c.interServerSecretKey,
	}

	_, _, err := c.post("/role/upgrade/admin", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientAuthentication) UpgradeToSupport(
	accountID int,
) error {
	body := httpAuthentication.UpgradeToSupportBody{
		AccountID:            accountID,
		InterServerSecretKey: c.interServerSecretKey,
	}

	_, _, err := c.post("/role/upgrade/support", body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientAuthentication) RecoveryPassword(
	twoFaCode,
	newPassword,
	accessToken string,
) error {
	body := httpAuthentication.RecoveryPasswordBody{
		TwoFaCode:   twoFaCode,
		NewPassword: newPassword,
	}

	_, _, err := c.postWithAccessToken("/password/recovery", accessToken, body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientAuthentication) ChangePassword(
	oldPassword,
	newPassword,
	accessToken string,
) error {
	body := httpAuthentication.ChangePasswordBody{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}

	_, _, err := c.postWithAccessToken("/password/change", accessToken, body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientAuthentication) post(path string, body any) ([]byte, *http.Response, error) {
	bytesBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(bytesBody))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, nil, errors.New("status code is not 200")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return responseBody, resp, nil
}

func (c *ClientAuthentication) postWithAccessToken(path, accessToken string, body any) ([]byte, *http.Response, error) {
	bytesBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(bytesBody))
	if err != nil {
		return nil, nil, err
	}
	req.AddCookie(&http.Cookie{Name: "Access-Token", Value: accessToken})
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, nil, errors.New("status code is not 200")
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return responseBody, resp, nil
}

func (c *ClientAuthentication) getWithAccessToken(path, accessToken string) ([]byte, error) {
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

	if resp.StatusCode >= 300 {
		return nil, errors.New("status code is not 200")
	}

	_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return _body, nil
}
