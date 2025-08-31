package unit

import (
	"forum-payment/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestControllerCreatePayment(t *testing.T) {
	testConfig.PrepareDB()

	paymentData, err := testConfig.paymentClient.CreatePayment(1, model.Frame, model.BTC, 10.4)
	assert.NoError(t, err)

	assert.Equal(t, 1, paymentData.PaymentID)
	assert.Equal(t, model.BTC, paymentData.Currency)

	paymentStatus, err := testConfig.paymentClient.StatusPayment(paymentData.PaymentID)
	assert.NoError(t, err)

	assert.Equal(t, model.Pending, paymentStatus.PaymentStatus)
}

func TestControllerPaidPayment(t *testing.T) {
	testConfig.PrepareDB()

	paymentData, err := testConfig.paymentClient.CreatePayment(1, model.Frame, model.BTC, 10.4)
	assert.NoError(t, err)

	err = testConfig.paymentClient.PaidPayment(paymentData.PaymentID)
	assert.NoError(t, err)

	_, err = testConfig.paymentClient.StatusPayment(paymentData.PaymentID)
	assert.NoError(t, err)
}

func TestControllerCancelPayment(t *testing.T) {
	testConfig.PrepareDB()

	paymentData, err := testConfig.paymentClient.CreatePayment(1, model.Frame, model.BTC, 10.4)
	assert.NoError(t, err)

	err = testConfig.paymentClient.CancelPayment(paymentData.PaymentID)
	assert.NoError(t, err)

	paymentStatus, err := testConfig.paymentClient.StatusPayment(paymentData.PaymentID)
	assert.NoError(t, err)

	assert.Equal(t, model.Canceled, paymentStatus.PaymentStatus)
}
