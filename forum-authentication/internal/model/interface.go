package model

import (
	"bytes"
	"context"
	"forum-authentication/pkg/api/forum-authorization"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type IAuthenticationService interface {
	Register(ctx context.Context, login, email, password string) (*AuthorizationData, error)
	Login(ctx context.Context, login, password, twoFaCode string) (*AuthorizationData, error)

	GenerateTwoFa(accountID int) (string, *bytes.Buffer, error)
	SetTwoFaKey(ctx context.Context, accountID int, twoFaKey string, twoFaCode string) error
	DeleteTwoFaKey(ctx context.Context, accountID int, twoFaCode string) error
	VerifyTwoFa(ctx context.Context, accountID int, twoFaCode string) (bool, error)
	UpgradeToAdmin(ctx context.Context, accountID int) error
	UpgradeToSupport(ctx context.Context, accountID int) error

	RecoveryPassword(ctx context.Context, accountID int, twoFaCode, newPassword string) error
	ChangePassword(ctx context.Context, accountID int, oldPassword, newPassword string) error
}

type IAuthenticationRepo interface {
	CreateAccount(ctx context.Context, login, email, password string) (int, error)

	SetTwoFaKey(ctx context.Context, accountID int, twoFaKey string) error
	DeleteTwoFaKey(ctx context.Context, accountID int) error
	UpgradeToAdmin(ctx context.Context, accountID int) error
	UpgradeToSupport(ctx context.Context, accountID int) error
	UpdatePassword(ctx context.Context, accountID int, newPassword string) error

	AccountByLogin(ctx context.Context, login string) ([]*Account, error)
	AccountByID(ctx context.Context, accountID int) ([]*Account, error)
}

type IAuthorizationClient interface {
	Authorization(ctx context.Context, accountID int) (*forum_authorization.JWTTokens, error)
	CheckAuthorization(ctx context.Context, request echo.Context) (*forum_authorization.AuthorizationData, error)
}

type IUserClient interface {
	CreateUser(ctx context.Context, accountID int, login string) error
}

type IAdminClient interface {
	CreateAdmin(ctx context.Context, accountID int) error
}

type IThreadClient interface {
	CreateAccountStatistic(ctx context.Context, accountID int) error
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
