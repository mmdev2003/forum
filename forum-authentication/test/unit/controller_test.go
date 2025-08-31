package unit

import (
	"forum-authentication/internal/model"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestControllerRegister(t *testing.T) {
	testConfig.PrepareDB()
	responseBody, resp, err := testConfig.authenticationClient.Register("user1", "user1@mail.com", "123456")
	assert.NoError(t, err)
	assert.Equal(t, 1, responseBody.AccountID)
	assert.Equal(t, 2, len(resp.Cookies()))

	for _, cookie := range resp.Cookies() {
		assert.Equal(t, true, cookie.Name == "Access-Token" || cookie.Name == "Refresh-Token")
	}
}

func TestControllerLogin(t *testing.T) {
	testConfig.PrepareDB()

	_, _, err := testConfig.authenticationClient.Register("user1", "user1@mail.com", "123456")
	assert.NoError(t, err)

	responseBody, resp, err := testConfig.authenticationClient.Login("user1", "123456", "")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Cookies()))

	for _, cookie := range resp.Cookies() {
		assert.Equal(t, true, cookie.Name == "Access-Token" || cookie.Name == "Refresh-Token")
	}

	assert.Equal(t, 1, responseBody.AccountID)
	assert.Equal(t, "user", responseBody.Role)
	assert.Equal(t, true, responseBody.IsTwoFaVerified)
}

func TestControllerGenerateTwoFa(t *testing.T) {
	testConfig.PrepareDB()
	_, _, err := testConfig.authenticationClient.Register("user1", "user1@mail.com", "123456")
	assert.NoError(t, err)

	qrCode, twoFaKey, err := testConfig.authenticationClient.GenerateTwoFa("user1")
	assert.NoError(t, err)

	assert.Equal(t, true, len(qrCode) > 0)
	assert.Equal(t, true, len(twoFaKey) > 0)
}

func TestControllerSetTwoFa(t *testing.T) {
	testConfig.PrepareDB()

	_, _, err := testConfig.authenticationClient.Register("user1", "user1@mail.com", "123456")
	assert.NoError(t, err)

	_, twoFaKey, err := testConfig.authenticationClient.GenerateTwoFa("user1")
	assert.NoError(t, err)

	twoFaCode, err := totp.GenerateCode(twoFaKey, time.Now())
	assert.NoError(t, err)

	err = testConfig.authenticationClient.SetTwoFa(twoFaKey, twoFaCode, "user1")
	assert.NoError(t, err)

	responseBody, _, err := testConfig.authenticationClient.Login("user1", "123456", "")
	assert.NoError(t, err)
	assert.Equal(t, false, responseBody.IsTwoFaVerified)

	verified, err := testConfig.authenticationClient.VerifyTwoFa(twoFaCode, "user1")
	assert.NoError(t, err)
	assert.Equal(t, true, verified.IsTwoFaVerified)
}

func TestControllerDeleteTwoFaKey(t *testing.T) {
	testConfig.PrepareDB()

	_, _, err := testConfig.authenticationClient.Register("user1", "user1@mail.com", "123456")
	assert.NoError(t, err)

	_, twoFaKey, err := testConfig.authenticationClient.GenerateTwoFa("user1")
	assert.NoError(t, err)

	twoFaCode, err := totp.GenerateCode(twoFaKey, time.Now())
	assert.NoError(t, err)

	err = testConfig.authenticationClient.SetTwoFa(twoFaKey, twoFaCode, "user1")
	assert.NoError(t, err)

	err = testConfig.authenticationClient.DeleteTwoFaKey(twoFaCode, "user1")
	assert.NoError(t, err)

	responseBody, _, err := testConfig.authenticationClient.Login("user1", "123456", "")
	assert.NoError(t, err)
	assert.Equal(t, true, responseBody.IsTwoFaVerified)
}

func TestControllerLoginWithTwoFa(t *testing.T) {
	testConfig.PrepareDB()

	_, _, err := testConfig.authenticationClient.Register("user1", "user1@mail.com", "123456")
	assert.NoError(t, err)

	_, twoFaKey, err := testConfig.authenticationClient.GenerateTwoFa("user1")
	assert.NoError(t, err)

	twoFaCode, err := totp.GenerateCode(twoFaKey, time.Now())
	assert.NoError(t, err)

	err = testConfig.authenticationClient.SetTwoFa(twoFaKey, twoFaCode, "user1")
	assert.NoError(t, err)

	responseBody, _, err := testConfig.authenticationClient.Login("user1", "123456", "")
	assert.NoError(t, err)

	assert.Equal(t, false, responseBody.IsTwoFaVerified)

	twoFaCode, err = totp.GenerateCode(twoFaKey, time.Now())
	assert.NoError(t, err)

	responseBody, _, err = testConfig.authenticationClient.Login("user1", "123456", twoFaCode)
	assert.NoError(t, err)

	assert.Equal(t, true, responseBody.IsTwoFaVerified)
}

func TestControllerUpgradeToAdmin(t *testing.T) {
	testConfig.PrepareDB()

	_, _, err := testConfig.authenticationClient.Register("user1", "user1@mail.com", "123456")
	assert.NoError(t, err)

	err = testConfig.authenticationClient.UpgradeToAdmin(1)

	responseBody, _, err := testConfig.authenticationClient.Login("user1", "123456", "")
	assert.NoError(t, err)

	assert.Equal(t, model.RoleAdmin, responseBody.Role)
}

func TestControllerChangePassword(t *testing.T) {
	testConfig.PrepareDB()

	_, _, err := testConfig.authenticationClient.Register("user1", "user1@mail.com", "123456")
	assert.NoError(t, err)

	err = testConfig.authenticationClient.ChangePassword("123456", "12345678", "user1")
	assert.NoError(t, err)

	_, _, err = testConfig.authenticationClient.Login("user1", "12345678", "")
	assert.NoError(t, err)
}

func TestControllerRecoveryPassword(t *testing.T) {
	testConfig.PrepareDB()

	_, _, err := testConfig.authenticationClient.Register("user1", "user1@mail.com", "123456")
	assert.NoError(t, err)

	_, twoFaKey, err := testConfig.authenticationClient.GenerateTwoFa("user1")
	assert.NoError(t, err)

	twoFaCode, err := totp.GenerateCode(twoFaKey, time.Now())
	assert.NoError(t, err)

	err = testConfig.authenticationClient.SetTwoFa(twoFaKey, twoFaCode, "user1")
	assert.NoError(t, err)

	err = testConfig.authenticationClient.RecoveryPassword(twoFaCode, "12345678", "user1")
	assert.NoError(t, err)

	_, _, err = testConfig.authenticationClient.Login("user1", "12345678", "")
	assert.NoError(t, err)
}

func TestControllerUpgradeToSupport(t *testing.T) {
	testConfig.PrepareDB()

	_, _, err := testConfig.authenticationClient.Register("user1", "user1@mail.com", "123456")
	assert.NoError(t, err)

	err = testConfig.authenticationClient.UpgradeToSupport(1)

	responseBody, _, err := testConfig.authenticationClient.Login("user1", "123456", "")
	assert.NoError(t, err)

	assert.Equal(t, model.RoleSupport, responseBody.Role)
}
