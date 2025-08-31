package unit

import (
	"forum-status/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestControllerCreatePaymentForStatus(t *testing.T) {
	testConfig.PrepareDB()

	_, err := testConfig.statusClient.CreatePaymentForStatus(1, 1, "btc", "user1")
	assert.NoError(t, err)

	accountStatus, err := testConfig.statusClient.GetStatusByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, model.Pending, accountStatus.Statuses[0].PaymentStatus)
	assert.Equal(t, 1, accountStatus.Statuses[0].StatusID)
	assert.Equal(t, 1, accountStatus.Statuses[0].AccountID)
}

func TestControllerConfirmPaymentForStatus(t *testing.T) {
	testConfig.PrepareDB()

	_, err := testConfig.statusClient.CreatePaymentForStatus(1, 1, "btc", "user1")
	assert.NoError(t, err)

	err = testConfig.statusClient.ConfirmPaymentForStatus(1, testConfig.interServerSecretKey)
	assert.NoError(t, err)

	accountStatus, err := testConfig.statusClient.GetStatusByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, model.Confirmed, accountStatus.Statuses[0].PaymentStatus)
}

func TestControllerAssignStatusToAccount(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.statusClient.AssignStatusToAccount(1, 1, testConfig.interServerSecretKey)
	assert.NoError(t, err)

	accountStatus, err := testConfig.statusClient.GetStatusByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, model.Confirmed, accountStatus.Statuses[0].PaymentStatus)
}

func TestControllerRevokeStatusFromAccount(t *testing.T) {
	testConfig.PrepareDB()

	err := testConfig.statusClient.AssignStatusToAccount(1, 1, testConfig.interServerSecretKey)
	assert.NoError(t, err)

	err = testConfig.statusClient.RevokeStatusFromAccount(1, 1, testConfig.interServerSecretKey)
	assert.NoError(t, err)

	accountStatus, err := testConfig.statusClient.GetStatusByAccountID(1)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(accountStatus.Statuses))
}

func TestControllerGetAllStatus(t *testing.T) {
	testConfig.PrepareDB()

	_, err := testConfig.statusClient.GetAllStatus()
	assert.NoError(t, err)
}
