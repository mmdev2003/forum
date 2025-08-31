package authorization

import (
	"context"
	"forum-authorization/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func New(db model.IDatabase) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

type AccountRepository struct {
	db model.IDatabase
}

func (accountRepo *AccountRepository) SetAccount(ctx context.Context,
	accountID int,
) error {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	_, err := accountRepo.db.Insert(ctx, SetAccount, args)
	if err != nil {
		return err
	}
	return nil
}

func (accountRepo *AccountRepository) UpdateRefreshToken(ctx context.Context,
	accountID int,
	refreshToken string,
) error {
	args := pgx.NamedArgs{
		"account_id":    accountID,
		"refresh_token": refreshToken,
	}
	err := accountRepo.db.Update(ctx, UpdateRefreshToken, args)
	if err != nil {
		return err
	}
	return nil
}

func (accountRepo *AccountRepository) AccountByID(ctx context.Context,
	accountID int,
) ([]*model.Account, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	rows, err := accountRepo.db.Select(ctx, AccountByID, args)
	if err != nil {
		return nil, err
	}

	var accounts []*model.Account
	err = pgxscan.ScanAll(&accounts, rows)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (accountRepo *AccountRepository) AccountByRefreshToken(ctx context.Context,
	refreshToken string,
) ([]*model.Account, error) {
	args := pgx.NamedArgs{
		"refresh_token": refreshToken,
	}
	rows, err := accountRepo.db.Select(ctx, AccountByRefreshToken, args)
	if err != nil {
		return nil, err
	}

	var accounts []*model.Account
	err = pgxscan.ScanAll(&accounts, rows)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
