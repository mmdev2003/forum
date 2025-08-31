package payment

import (
	"context"
	"encoding/json"
	"forum-payment/internal/model"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

func New(
	paymentRepo model.IPaymentRepo,
	frameClient model.IFrameClient,
	statusClient model.IStatusClient,
	btcAddress string,
) *ServicePayment {
	return &ServicePayment{
		paymentRepo:  paymentRepo,
		frameClient:  frameClient,
		statusClient: statusClient,
		btcAddress:   btcAddress,
	}
}

type ServicePayment struct {
	paymentRepo  model.IPaymentRepo
	frameClient  model.IFrameClient
	statusClient model.IStatusClient
	btcAddress   string
}

func (paymentService *ServicePayment) CreatePayment(ctx context.Context,
	accountID int,
	productType model.ProductType,
	currency model.Currency,
	amountUSD float32,
) (string, string, int, error) {
	var amount string
	var address string
	var err error

	if currency == model.BTC {
		address = paymentService.btcAddress
		amount, err = paymentService.generateBTCPrice(ctx,
			amountUSD,
			address,
		)
		if err != nil {
			return "", "", 0, err
		}
	}

	paymentID, err := paymentService.paymentRepo.CreatePayment(
		ctx,
		accountID,
		productType,
		currency,
		amount,
		address,
	)
	if err != nil {
		return "", "", 0, err
	}

	return amount, address, paymentID, nil
}

func (paymentService *ServicePayment) PaidPayment(ctx context.Context,
	paymentID int,
) error {
	err := paymentService.paymentRepo.PaidPayment(ctx, paymentID)
	if err != nil {
		return err
	}
	return nil
}

func (paymentService *ServicePayment) CancelPayment(ctx context.Context,
	paymentID int,
) error {
	err := paymentService.paymentRepo.CancelPayment(ctx, paymentID)
	if err != nil {
		return err
	}
	return nil
}

func (paymentService *ServicePayment) StatusPayment(ctx context.Context,
	paymentID int,
) (model.PaymentStatus, error) {
	payments, err := paymentService.paymentRepo.PaymentByID(
		ctx,
		paymentID,
	)
	if err != nil {
		return "", err
	}
	if len(payments) == 0 {
		return "", model.ErrPaymentNotFound
	}
	return payments[0].Status, nil
}

func (paymentService *ServicePayment) ConfirmPayment(ctx context.Context,
	paymentID int,
	txID string,
	productType model.ProductType,
) error {
	switch productType {
	case model.Frame:
		err := paymentService.frameClient.ConfirmPaymentForFrame(ctx, paymentID)
		if err != nil {
			return err
		}
	case model.Status:
		err := paymentService.statusClient.ConfirmPaymentForStatus(ctx, paymentID)
		if err != nil {
			return err
		}
	}

	err := paymentService.paymentRepo.ConfirmPayment(ctx, paymentID, txID)
	if err != nil {
		return err
	}
	return nil
}

func (paymentService *ServicePayment) generateBTCPrice(ctx context.Context,
	amountUSD float32,
	address string,
) (string, error) {
	countRandomLastDigit := 3
	btcPrice, err := getBtcPrice()
	if err != nil {
		return "", err
	}
	amountInBTC := amountUSD / btcPrice
	amountInBTCStr := strconv.FormatFloat(float64(amountInBTC), 'f', 28, 32)

	randomLastDigit := generateRandomDigits(countRandomLastDigit)
	amountInBTCStr = amountInBTCStr[:len(amountInBTCStr)-countRandomLastDigit] + randomLastDigit

	payments, err := paymentService.paymentRepo.PaymentByAddressAndAmount(ctx, address, amountInBTCStr)
	if err != nil {
		return "", err
	}

	if len(payments) > 0 {
		for len(payments) > 0 {
			randomLastDigit = generateRandomDigits(countRandomLastDigit)
			amountInBTCStr = amountInBTCStr[:len(amountInBTCStr)-countRandomLastDigit] + randomLastDigit

			payments, err = paymentService.paymentRepo.PaymentByAddressAndAmount(ctx, address, amountInBTCStr)
			if err != nil {
				return "", err
			}
		}
	}

	return amountInBTCStr, nil
}

func getBtcPrice() (float32, error) {
	type PriceResponse struct {
		BTC struct {
			USD float32 `json:"usd"`
		} `json:"bitcoin"`
	}
	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var priceData PriceResponse
	err = json.Unmarshal(body, &priceData)
	if err != nil {
		return 0, err
	}
	return priceData.BTC.USD, nil
}

func generateRandomDigits(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte('0' + byte(rand.Intn(10)))
	}
	return sb.String()
}
