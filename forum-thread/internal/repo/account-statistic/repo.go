package account_statistic

import (
	"context"
	"forum-thread/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func New(
	db model.IDatabase,
) *RepoAccountStatistic {
	return &RepoAccountStatistic{
		db: db,
	}
}

type RepoAccountStatistic struct {
	db model.IDatabase
}

func (r *RepoAccountStatistic) CreateAccountStatistic(ctx context.Context,
	accountID int,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	accountStatisticID, err := r.db.Insert(ctx, CreateAccountStatistic, args)
	if err != nil {
		return 0, err
	}

	return accountStatisticID, nil
}
func (r *RepoAccountStatistic) AddSentMessagesToTopicsCount(ctx context.Context,
	accountID int,
) error {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	err := r.db.Update(ctx, AddSentMessagesToTopicsCount, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoAccountStatistic) AddCreatedTopicsCount(ctx context.Context,
	accountID int,
) error {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	err := r.db.Update(ctx, AddCreatedTopicsCount, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoAccountStatistic) StatisticByAccountID(ctx context.Context,
	accountID int,
) ([]*model.AccountStatistic, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	rows, err := r.db.Select(ctx, StatisticByAccountID, args)
	if err != nil {
		return nil, err
	}

	var accountStatistic []*model.AccountStatistic
	err = pgxscan.ScanAll(&accountStatistic, rows)
	if err != nil {
		return nil, err
	}

	return accountStatistic, nil
}
