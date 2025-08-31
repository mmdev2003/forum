package model

import (
	"context"
	"forum-status/pkg/api/forum-authorization"
	"forum-status/pkg/api/forum-payment"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"time"
)

type IStatusService interface {
	CreatePaymentForStatus(ctx context.Context, statusID, accountID, duration int, currency string) (*forum_payment.PaymentData, error)
	AssignStatusToAccount(ctx context.Context, statusID, accountID int) error
	ConfirmPaymentForStatus(ctx context.Context, paymentID int) error
	RevokeStatus(ctx context.Context, statusID, accountID int) error
	StatusesByAccountID(ctx context.Context, accountID int) ([]*Status, error)
}

type IStatusRepo interface {
	CreateStatus(ctx context.Context, statusID, accountID, paymentID int, expirationAt time.Time) (int, error)
	ConfirmPaymentForStatus(ctx context.Context, paymentID int) error
	DeleteStatus(ctx context.Context, statusID, accountID int) error
	StatusesByAccountID(ctx context.Context, accountID int) ([]*Status, error)
}

type IAuthorizationClient interface {
	CheckAuthorization(ctx context.Context, request echo.Context) (*forum_authorization.AuthorizationData, error)
}

type IPaymentClient interface {
	CreatePayment(ctx context.Context, accountID int, productType string, currency string, amountUSD float32) (*forum_payment.PaymentData, error)
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
