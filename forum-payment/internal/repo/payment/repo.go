package payment

import (
	"context"
	"forum-payment/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func New(db model.IDatabase) *RepoPayment {
	return &RepoPayment{
		db: db,
	}
}

type RepoPayment struct {
	db model.IDatabase
}

func (paymentRepo *RepoPayment) CreatePayment(ctx context.Context,
	accountID int,
	productType model.ProductType,
	currency model.Currency,
	amount string,
	address string,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id":   accountID,
		"product_type": productType,
		"currency":     currency,
		"amount":       amount,
		"address":      address,
		"status":       model.Pending,
	}
	paymentID, err := paymentRepo.db.Insert(ctx, CreatePayment, args)
	if err != nil {
		return 0, err
	}
	return paymentID, nil
}

func (paymentRepo *RepoPayment) PaidPayment(ctx context.Context,
	paymentID int,
) error {
	args := pgx.NamedArgs{
		"payment_id": paymentID,
	}
	err := paymentRepo.db.Update(ctx, PaidPayment, args)
	if err != nil {
		return err
	}
	return nil
}

func (paymentRepo *RepoPayment) CancelPayment(ctx context.Context,
	paymentID int,
) error {
	args := pgx.NamedArgs{
		"payment_id": paymentID,
		"status":     model.Canceled,
	}
	err := paymentRepo.db.Update(ctx, CancelPayment, args)
	if err != nil {
		return err
	}
	return nil
}

func (paymentRepo *RepoPayment) ConfirmPayment(ctx context.Context,
	paymentID int,
	txID string,
) error {
	args := pgx.NamedArgs{
		"payment_id": paymentID,
		"tx_id":      txID,
		"status":     model.Confirmed,
	}
	err := paymentRepo.db.Update(ctx, ConfirmPayment, args)
	if err != nil {
		return err
	}
	return nil
}

func (paymentRepo *RepoPayment) PaymentByID(ctx context.Context,
	paymentID int,
) ([]*model.Payment, error) {
	args := pgx.NamedArgs{
		"payment_id": paymentID,
	}
	rows, err := paymentRepo.db.Select(ctx, PaymentByID, args)
	if err != nil {
		return nil, err
	}

	var payment []*model.Payment
	err = pgxscan.ScanAll(&payment, rows)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (paymentRepo *RepoPayment) PaymentByTxID(ctx context.Context,
	txID string,
) ([]*model.Payment, error) {
	args := pgx.NamedArgs{
		"tx_id": txID,
	}
	rows, err := paymentRepo.db.Select(ctx, PaymentByTxID, args)
	if err != nil {
		return nil, err
	}

	var payment []*model.Payment
	err = pgxscan.ScanAll(&payment, rows)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (paymentRepo *RepoPayment) PendingPayments(ctx context.Context) ([]*model.Payment, error) {
	args := pgx.NamedArgs{
		"status": model.Pending,
	}
	rows, err := paymentRepo.db.Select(ctx, PendingPayments, args)
	if err != nil {
		return nil, err
	}

	var payments []*model.Payment
	err = pgxscan.ScanAll(&payments, rows)
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (paymentRepo *RepoPayment) PaymentByAddressAndAmount(ctx context.Context,
	address string,
	amount string,
) ([]*model.Payment, error) {
	args := pgx.NamedArgs{
		"address": address,
		"amount":  amount,
	}
	rows, err := paymentRepo.db.Select(ctx, PaymentByAddressAndAmount, args)
	if err != nil {
		return nil, err
	}

	var payment []*model.Payment
	err = pgxscan.ScanAll(&payment, rows)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
