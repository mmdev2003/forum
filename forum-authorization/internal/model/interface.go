package model

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type IAuthorizationService interface {
	CreateTokens(ctx context.Context, accountID int, role string, twoFaStatus bool) (*JWTTokens, error)
	CheckToken(token string) (*TokenPayload, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*JWTTokens, error)
}

type IAuthorizationRepo interface {
	SetAccount(ctx context.Context, accountID int) error
	UpdateRefreshToken(ctx context.Context, accountID int, refreshToken string) error
	AccountByID(ctx context.Context, accountID int) ([]*Account, error)
	AccountByRefreshToken(ctx context.Context, refreshToken string) ([]*Account, error)
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
