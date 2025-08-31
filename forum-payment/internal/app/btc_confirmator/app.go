package btc_confirmator

import (
	"context"
	"encoding/json"
	"fmt"
	"forum-payment/internal/model"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func Run(
	paymentService model.IPaymentService,
	paymentRepo model.IPaymentRepo,
	btcAddress string,
) {
	for {
		ctx := context.Background()
		pendingPayments, err := paymentRepo.PendingPayments(ctx)
		if err != nil {
			slog.Error(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		fmt.Println(pendingPayments)

		receivedBTC, txIDs, err := getReceivedBTC(ctx, btcAddress, paymentRepo)
		if err != nil {
			slog.Error(err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		fmt.Println(receivedBTC)

		if len(pendingPayments) > 0 && len(receivedBTC) > 0 {
			for _, payment := range pendingPayments {
				for i, amount := range receivedBTC {
					if payment.Amount == fmt.Sprintf("%.8f", amount) {
						txID := txIDs[i]
						err = paymentService.ConfirmPayment(ctx, payment.ID, txID, payment.ProductType)
						if err != nil {
							slog.Error(err.Error())
						}
					} else {

					}
				}
			}
		}

		time.Sleep(60 * time.Second)
	}
}

func getReceivedBTC(
	ctx context.Context,
	btcAddress string,
	paymentRepo model.IPaymentRepo,
) ([]float32, []string, error) {
	url := fmt.Sprintf("https://btcscan.org/api/address/%s/txs", btcAddress)

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var txs []struct {
		Status struct {
			Confirmed bool `json:"confirmed"`
		} `json:"status"`
		TxID string `json:"txid"`
		Vout []struct {
			Value               float64 `json:"value"`
			ScriptPubKeyAddress string  `json:"scriptpubkey_address"`
		} `json:"vout"`
	}

	err = json.Unmarshal(body, &txs)
	if err != nil {
		return nil, nil, err
	}

	var receivedBTC []float32
	var txIDs []string

	for _, tx := range txs {
		if tx.Status.Confirmed {
			for _, vout := range tx.Vout {
				if vout.ScriptPubKeyAddress == btcAddress {
					txID := tx.TxID
					payment, err := paymentRepo.PaymentByTxID(ctx, txID)
					if err != nil || payment != nil {
						continue
					}

					value := float32(vout.Value / 100000000)
					receivedBTC = append(receivedBTC, value)
					txIDs = append(txIDs, txID)
				}
			}
		}
	}

	return receivedBTC, txIDs, nil
}
