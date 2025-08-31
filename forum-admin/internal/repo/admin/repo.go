package admin

import (
	"context"
	"forum-admin/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func New(
	db model.IDatabase,
) *RepoAdmin {
	return &RepoAdmin{db}
}

type RepoAdmin struct {
	db model.IDatabase
}

func (r *RepoAdmin) CreateAdmin(ctx context.Context,
	accountID int,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	adminID, err := r.db.Insert(ctx, CreateAdmin, args)
	if err != nil {
		return 0, err
	}

	return adminID, nil
}

func (r *RepoAdmin) AllAdmin(ctx context.Context,
) ([]*model.Admin, error) {
	rows, err := r.db.Select(ctx, AllAdmin, pgx.NamedArgs{})
	if err != nil {
		return nil, err
	}

	var admins []*model.Admin
	err = pgxscan.ScanAll(&admins, rows)
	if err != nil {
		return nil, err
	}

	return admins, nil
}
