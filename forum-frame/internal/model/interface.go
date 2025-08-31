package model

import (
	"context"
	"forum-frame/pkg/api/forum-authorization"
	"forum-frame/pkg/api/forum-payment"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"time"
)

type IFrameService interface {
	CreatePaymentForFrame(ctx context.Context, frameID, accountID, duration int, currency string) (*forum_payment.PaymentData, error)
	AddNewFrame(ctx context.Context, frameFile []byte, monthlyPrice, foreverPrice float32, name string) error

	ConfirmPaymentForFrame(ctx context.Context, paymentID int) error
	ChangeCurrentFrame(ctx context.Context, dbFrameID, accountID int) error

	FramesByAccountID(ctx context.Context, accountID int) ([]*Frame, *CurrentFrame, error)
	AllFrame(ctx context.Context) ([]*FrameData, error)
	DownloadFrame(ctx context.Context, frameID int) ([]byte, error)
}

type IFrameRepo interface {
	CreateFrame(ctx context.Context, frameID, accountID, paymentID int, expirationAt time.Time) (int, error)
	CreateFrameData(ctx context.Context, monthlyPrice, foreverPrice float32, name, fileID string) error
	ConfirmPaymentForFrame(ctx context.Context, paymentID int) error

	UploadFrame(frameFile []byte, name string) (string, error)
	ChangeCurrentFrame(ctx context.Context, dbFrameID, accountID int) error

	DownloadFrame(ctx context.Context, frameID int) ([]byte, error)
	FramesByAccountID(ctx context.Context, accountID int) ([]*Frame, error)
	AllFrame(ctx context.Context) ([]*FrameData, error)
	FrameDataByID(ctx context.Context, frameID int) ([]*FrameData, error)
	CurrentFrameByAccountID(ctx context.Context, accountID int) ([]*CurrentFrame, error)
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

type IWeedFS interface {
	Upload(file []byte, filename string, size int64, collection string) (string, error)
	Download(FileID string) ([]byte, error)
	Update(file []byte, fileID string, filename string, size int64, collection string) error
}
