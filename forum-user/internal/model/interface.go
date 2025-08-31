package model

import (
	"context"
	"forum-user/pkg/api/forum-authorization"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
)

type IUserService interface {
	CreateUser(ctx context.Context, accountID int, login string) (int, error)
	BanUser(ctx context.Context, fromAccountID, toAccountID int) error
	NewWarningFromAdmin(ctx context.Context, adminAccountID, toAccountID int, warningType, adminLogin string) error

	UploadAvatar(ctx context.Context, accountID int, avatar []byte) error
	UploadHeader(ctx context.Context, accountID int, header []byte) error

	UnbanUser(ctx context.Context, fromAccountID, toAccountID int) error

	DownloadAvatar(ctx context.Context, fileID string) ([]byte, error)
	DownloadHeader(ctx context.Context, fileID string) ([]byte, error)
	UserByAccountID(ctx context.Context, accountID int) (*User, error)
	UserByLogin(ctx context.Context, login string) ([]*User, error)
	UsersByLoginSearch(ctx context.Context, login string) ([]*UserSearch, error)
	BanByAccountID(ctx context.Context, toAccountID int) ([]*UserBan, error)
	AllWarningFromAdmin(ctx context.Context, toAccountID int) ([]*WarningFromAdmin, error)
}

type IUserRepo interface {
	CreateUser(ctx context.Context, accountID int, login string) (int, error)
	CreateUserSearchIndex(ctx context.Context, accountID int, login string) error
	CreateUserBan(ctx context.Context, fromAccountID, toAccountID int) error
	CreateWarningFromAdmin(ctx context.Context, adminAccountID, toAccountID int, warningType, warningText, adminLogin string) error

	UploadAvatar(ctx context.Context, accountID int, avatar []byte) (string, error)
	UpdateAvatarUrl(ctx context.Context, accountID int, avatarUrl string) error
	UploadHeader(ctx context.Context, accountID int, header []byte) (string, error)
	UpdateHeaderUrl(ctx context.Context, accountID int, headerUrl string) error

	DeleteUserBan(ctx context.Context, fromAccountID, toAccountID int) error

	DownloadAvatar(ctx context.Context, fileID string) ([]byte, error)
	DownloadHeader(ctx context.Context, fileID string) ([]byte, error)
	UserByAccountID(ctx context.Context, accountID int) ([]*User, error)
	UsersByLoginSearch(ctx context.Context, login string) ([]*UserSearch, error)
	UserByLogin(ctx context.Context, login string) ([]*User, error)
	BanByAccountID(ctx context.Context, toAccountID int) ([]*UserBan, error)
	AllWarningFromAdmin(ctx context.Context, toAccountID int) ([]*WarningFromAdmin, error)
}

type IAuthorizationClient interface {
	CheckAuthorization(ctx context.Context, request echo.Context) (*forum_authorization.AuthorizationData, error)
}

type INotificationClient interface {
	NewWarningFromAdminNotification(ctx context.Context,
		toAccountID,
		adminAccountID int,
		warningText,
		adminLogin string,
	) error
}
type IFullTextSearchEngine interface {
	CreateIndexes(indexes []string) error
	AddDocuments(indexName string, documents []any) error
	SimpleSearch(indexName, query string) ([]byte, error)
}

type IDatabase interface {
	Insert(ctx context.Context, query string, args ...any) (int, error)
	Select(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Delete(ctx context.Context, query string, args ...any) error
	Update(ctx context.Context, query string, args ...any) error

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

type IMessageBroker interface {
	Publish(ctx context.Context, queueName string, body any) error
	Subscribe(queueName string, handler func(event amqp091.Delivery) error)
}
