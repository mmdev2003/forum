package model

import (
	"context"
	forumauthorization "forum-notification/pkg/api/forum-authorization"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type INotificationFilter interface {
	GetFilters(ctx context.Context, accountID int) ([]NotificationType, error)
	UpsertFilters(ctx context.Context, accountID int, filters []NotificationType) error
	IsNotificationEnabled(ctx context.Context, accountID int, notificationType NotificationType) (bool, error)
}

type INotificationFilterRepo interface {
	GetFilters(ctx context.Context, accountID int) ([]NotificationType, error)
	UpsertFilters(ctx context.Context, accountID int, filters []NotificationType) error
	IsNotificationEnabled(ctx context.Context, accountID int, notificationType NotificationType) (bool, error)
}

type INotification interface {
	CreateMessageFromTopicNotification(ctx context.Context,
		accountID,
		messageID,
		replierAccountID,
		topicID int,
		messageText,
		topicName,
		replierLogin string,
	) (int, error)
	CreateMessageReplyFromTopicNotification(ctx context.Context,
		accountID,
		messageID,
		replierAccountID,
		topicID int,
		messageText,
		topicName,
		replierLogin string,
	) (int, error)
	CreateLikeMessageFromTopicNotification(ctx context.Context,
		accountID,
		messageID,
		likerAccountID,
		topicID int,
		messageText,
		topicName,
		likerLogin string,
	) (int, error)
	CreateTopicClosedNotification(ctx context.Context,
		accountID,
		adminAccountID,
		topicID int,
		topicName,
		adminLogin string,
	) (int, error)
	CreateResponseToSupportRequestNotification(ctx context.Context,
		accountID,
		supportRequestID int,
	) (int, error)
	CreateStatusReceivedNotification(ctx context.Context,
		accountID int,
		statusName string,
	) (int, error)
	CreateFrameReceivedNotification(ctx context.Context,
		accountID int,
		frameName string,
	) (int, error)
	CreateMessageFromDialogNotification(ctx context.Context,
		accountID,
		messageID,
		dialogID,
		senderAccountID int,
		messageText,
		senderLogin string,
	) (int, error)
	CreateMentionFromTopicNotification(ctx context.Context,
		accountID,
		messageID,
		mentionAccountID int,
		messageText,
		topicName,
		mentionLogin string,
	) (int, error)
	CreateWarningFromAdminNotification(ctx context.Context,
		accountID,
		adminAccountID int,
		warningText,
		adminLogin string,
	) (int, error)
	GetNotificationsByAccountID(ctx context.Context, id int) ([]interface{}, error)
}

type INotificationRepo interface {
	CreateMessageFromTopicNotification(ctx context.Context,
		accountID,
		messageID,
		replierAccountID,
		topicID int,
		messageText,
		topicName,
		replierLogin string,
	) (int, error)
	CreateMessageReplyFromTopicNotification(ctx context.Context,
		accountID,
		messageID,
		replierAccountID,
		topicID int,
		messageText,
		topicName,
		replierLogin string,
	) (int, error)
	CreateLikeMessageFromTopicNotification(ctx context.Context,
		accountID,
		messageID,
		likerAccountID,
		topicID int,
		messageText,
		topicName,
		likerLogin string,
	) (int, error)
	CreateTopicClosedNotification(ctx context.Context,
		accountID,
		adminAccountID,
		topicID int,
		topicName,
		adminLogin string,
	) (int, error)
	CreateResponseToSupportRequestNotification(ctx context.Context,
		accountID,
		supportRequestID int,
	) (int, error)
	CreateStatusReceivedNotification(ctx context.Context,
		accountID int,
		statusName string,
	) (int, error)
	CreateFrameReceivedNotification(ctx context.Context,
		accountID int,
		frameName string,
	) (int, error)
	CreateMessageFromDialogNotification(ctx context.Context,
		accountID,
		messageID,
		dialogID,
		senderAccountID int,
		messageText,
		senderLogin string,
	) (int, error)
	CreateMentionFromTopicNotification(ctx context.Context,
		accountID,
		messageID,
		mentionAccountID int,
		messageText,
		topicName,
		mentionLogin string,
	) (int, error)
	CreateWarningFromAdminNotification(ctx context.Context,
		accountID,
		adminAccountID int,
		warningText,
		adminLogin string,
	) (int, error)
	GetNotificationsByAccountID(ctx context.Context, id int) ([]interface{}, error)
}

type IWsConnManager interface {
	AddConnection(accountID int, conn *websocket.Conn)
	RemoveConnection(accountID int)
	SendMsg(toAccountID int, msg []byte) error
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

type IAuthorizationClient interface {
	CheckAuthorization(ctx context.Context, request echo.Context) (*forumauthorization.AuthorizationData, error)
}
