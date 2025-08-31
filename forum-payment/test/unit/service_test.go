package unit

import (
	"context"
	"forum-payment/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceConfirmFramePayment(t *testing.T) {
	testConfig.PrepareDB()
	ctx := context.Background()

	_, _, _, err := testConfig.paymentService.CreatePayment(ctx,
		1,
		model.Frame,
		model.BTC,
		10.4,
	)
	assert.NoError(t, err)

	err = testConfig.paymentService.ConfirmPayment(ctx, 1, "23423", model.Frame)
	assert.NoError(t, err)

	payment, err := testConfig.paymentRepo.PaymentByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, model.Confirmed, payment[0].Status)
	assert.Equal(t, "23423", payment[0].TxID)
}

func TestServiceConfirmStatusPayment(t *testing.T) {
	testConfig.PrepareDB()
	ctx := context.Background()

	_, _, _, err := testConfig.paymentService.CreatePayment(ctx,
		1,
		model.Status,
		model.BTC,
		10.4,
	)
	assert.NoError(t, err)

	err = testConfig.paymentService.ConfirmPayment(ctx, 1, "23423", model.Status)
	assert.NoError(t, err)

	payment, err := testConfig.paymentRepo.PaymentByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, model.Confirmed, payment[0].Status)
	assert.Equal(t, "23423", payment[0].TxID)
}
