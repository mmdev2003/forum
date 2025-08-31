package model

import (
	"context"
	"forum-dialog/pkg/api/forum-authorization"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

type IDialogService interface {
	CreateDialog(ctx context.Context, account1ID, account2ID int) (int, error)
	AddMessageToDialog(ctx context.Context, dialogID, fromAccountID, toAccountID int, text string, filesURLS []string) (int, error)

	MarkDialogAsStarred(ctx context.Context, dialogID, accountID int) error
	MarkMessagesAsRead(ctx context.Context, dialogID int) error
	UploadFile(ctx context.Context, file []byte, fullName string) (string, error)

	DeleteDialog(ctx context.Context, dialogID int) error

	DialogsByAccountID(ctx context.Context, accountID int) ([]*Dialog, error)
	MessagesByDialogID(ctx context.Context, dialogID int) ([]*Message, []*File, error)
	DownloadFile(ctx context.Context, fileURL string) ([]byte, error)
}

type IDialogRepo interface {
	CreateDialog(ctx context.Context, account1ID, account2ID int) (int, error)
	CreateMessage(ctx context.Context, dialogID, fromAccountID, toAccountID int, text string) (int, error)
	CreateFile(ctx context.Context, size int, fileURL, name, extension string) error

	MarkDialogAsStarred(ctx context.Context, dialogID, accountID int) error
	UpdateLastMessageAt(ctx context.Context, dialogID int) error
	MarkMessagesAsRead(ctx context.Context, dialogID int) error
	UploadFile(ctx context.Context, file []byte, fullName string) (string, error)
	AddFileToMessage(ctx context.Context, messageID int, fileURL string) error

	DeleteDialog(ctx context.Context, dialogID int) error

	DialogsByAccountID(ctx context.Context, accountID int) ([]*Dialog, error)
	MessagesByDialogID(ctx context.Context, dialogID int) ([]*Message, error)
	FilesByMessageID(ctx context.Context, messageID int) ([]*File, error)
	DownloadFile(ctx context.Context, fileURL string) ([]byte, error)

	CtxWithTx(ctx context.Context) (context.Context, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context)
}

type IAuthorizationClient interface {
	CheckAuthorization(ctx context.Context, request echo.Context) (*forum_authorization.AuthorizationData, error)
}

type IWsConnManager interface {
	AddConnection(ctx context.Context, accountID int, conn *websocket.Conn) error
	RemoveConnection(accountID int)
	SendMsg(ctx context.Context, toAccountID int, msg []byte) error
}

type IRedis interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
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

type IMessageBroker interface {
	Publish(ctx context.Context, queueName string, body any) error
	Subscribe(queueName string, handler func(event amqp091.Delivery) error)
}

type IWeedFS interface {
	Upload(file []byte, filename string, size int64, collection string) (string, error)
	Download(FileID string) ([]byte, error)
	Update(file []byte, fileID string, filename string, size int64, collection string) error
}
