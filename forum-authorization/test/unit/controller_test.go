package unit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestControllerAuthorization(t *testing.T) {
	testConfig.PrepareDB()

	jwtTokens, err := testConfig.authorizationClient.Authorization(1, "user", false)
	assert.NoError(t, err)

	authData, err := testConfig.authorizationClient.CheckAuthorization(jwtTokens.AccessToken)
	assert.NoError(t, err)

	assert.Equal(t, 1, authData.AccountID)
	assert.Equal(t, "user", authData.Role)
	assert.Equal(t, false, authData.TwoFaStatus)
}

func TestControllerRefreshTokens(t *testing.T) {
	testConfig.PrepareDB()

	jwtTokens, err := testConfig.authorizationClient.Authorization(1, "user", false)
	assert.NoError(t, err)

	resp, err := testConfig.authorizationClient.RefreshTokens(jwtTokens.RefreshToken)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp.Cookies()))

	for _, cookie := range resp.Cookies() {
		assert.Equal(t, true, cookie.Name == "Access-Token" || cookie.Name == "Refresh-Token")
	}
}
