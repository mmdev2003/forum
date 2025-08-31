package model

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type IAdminService interface {
	CreateAdmin(ctx context.Context, accountID int) (int, error)
	AllAdmin(ctx context.Context) ([]*Admin, error)
}

type IAdminRepo interface {
	CreateAdmin(ctx context.Context, accountID int) (int, error)
	AllAdmin(ctx context.Context) ([]*Admin, error)
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
