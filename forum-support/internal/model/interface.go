package model

import (
	"context"
	forum_authorization "forum-support/pkg/api/forum-authorization"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
)

type ISupportRequestService interface {
	CreateRequest(ctx context.Context, accountID int, title, description string) (int, error)
	OpenRequest(ctx context.Context, supportRequestID int, role string) error
	CloseRequest(ctx context.Context, supportRequestID, actorAccountID int, role string) error
	GetRequestById(ctx context.Context, supportRequestID int) (*SupportRequest, error)
	GetRequests(ctx context.Context) ([]SupportRequest, error)
	GetRequestsWithStatus(ctx context.Context, status RequestStatus) ([]SupportRequest, error)
}

type ISupportRequestRepo interface {
	CreateRequest(ctx context.Context, accountID int, title, description string) (int, error)
	OpenRequest(ctx context.Context, supportRequestID int) error
	CloseRequest(ctx context.Context, supportRequestID int) error
	GetRequestById(ctx context.Context, supportRequestID int) (*SupportRequest, error)
	GetRequests(ctx context.Context) ([]SupportRequest, error)
	GetRequestsWithStatus(ctx context.Context, status RequestStatus) ([]SupportRequest, error)
}

type IDialogService interface {
	CreateDialog(ctx context.Context, supportRequestID, userAccountID int) (int, error)

	AddMessageToDialog(ctx context.Context, dialogID, fromAccountID, toAccountID int, text string) (int, error)

	MarkMessagesAsRead(ctx context.Context, dialogID int, accountID int, role string) error

	AllDialogs(ctx context.Context, accountID int, role string) ([]*Dialog, error)
	MessagesByDialogID(ctx context.Context, dialogID int, accountID int, role string) ([]*Message, error)
}

type IDialogRepo interface {
	CreateDialog(ctx context.Context, supportRequestID, userAccountID int) (int, error)

	AddMessageToDialog(ctx context.Context, dialogID, fromAccountID, toAccountID int, text string) (int, error)

	MarkMessagesAsRead(ctx context.Context, dialogID, requesterID int, as_support bool) error

	AllDialogs(ctx context.Context) ([]*Dialog, error)
	GetDialogByID(ctx context.Context, dialogID int) (*Dialog, error)
	DialogsByAccountID(ctx context.Context, accountID int) ([]*Dialog, error)
	MessagesByDialogID(ctx context.Context, dialogID int) ([]*Message, error)
}

type IAuthorizationClient interface {
	CheckAuthorization(ctx context.Context, request echo.Context) (*forum_authorization.AuthorizationData, error)
}

type INotificationClient interface {
	SendSupportRequestClosedNotification(requesterAccountID, supportRequestID int) error
	SendStatusReceivedNotificationRequest(receiverAccountID int, statusName RequestStatus) error
	SendResponseToSupportRequestNotificationRequest(requesterAccountID, supportRequestID int) error
}

type IWsConnManager interface {
	AddConnection(ctx context.Context, accountID int, conn *websocket.Conn) error
	RemoveConnection(accountID int)
	SendMsg(ctx context.Context, toAccountID int, msg []byte) error
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

type IRedis interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type IMessageBroker interface {
	Publish(ctx context.Context, queueName string, body any) error
	Subscribe(queueName string, handler func(event amqp091.Delivery) error)
}
