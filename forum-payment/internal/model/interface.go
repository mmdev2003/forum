package model

import (
	"context"
	"forum-payment/pkg/api/forum-authorization"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type IPaymentService interface {
	CreatePayment(ctx context.Context,
		accountID int,
		productType ProductType,
		currency Currency,
		amountUSD float32,
	) (string, string, int, error)

	PaidPayment(ctx context.Context, paymentID int) error
	CancelPayment(ctx context.Context, paymentID int) error
	ConfirmPayment(ctx context.Context, paymentID int, txID string, productType ProductType) error

	StatusPayment(ctx context.Context, paymentID int) (PaymentStatus, error)
}

type IPaymentRepo interface {
	CreatePayment(ctx context.Context,
		accountID int,
		productType ProductType,
		currency Currency,
		amount string,
		address string,
	) (int, error)

	PaidPayment(ctx context.Context, paymentID int) error
	CancelPayment(ctx context.Context, paymentID int) error
	ConfirmPayment(ctx context.Context, paymentID int, txID string) error

	PaymentByID(ctx context.Context, paymentID int) ([]*Payment, error)
	PaymentByTxID(ctx context.Context, txID string) ([]*Payment, error)
	PendingPayments(ctx context.Context) ([]*Payment, error)
	PaymentByAddressAndAmount(ctx context.Context, address string, amount string) ([]*Payment, error)
}

type IAuthorizationClient interface {
	CheckAuthorization(ctx context.Context, request echo.Context) (*forum_authorization.AuthorizationData, error)
}

type IDatabase interface {
	Insert(ctx context.Context, query string, args ...interface{}) (int, error)
	Select(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Delete(ctx context.Context, query string, args ...interface{}) error
	Update(ctx context.Context, query string, args ...interface{}) error

	CtxWithTx(ctx context.Context) (context.Context, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context)

	CreateTable(ctx context.Context, query string) error
	DropTable(ctx context.Context, query string) error
}

type IFrameClient interface {
	ConfirmPaymentForFrame(ctx context.Context, paymentID int) error
}

type IStatusClient interface {
	ConfirmPaymentForStatus(ctx context.Context, paymentID int) error
}
