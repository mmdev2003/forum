package authentication

import (
	"context"
	"forum-authentication/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"time"
)

func New(db model.IDatabase) *RepoAuthentication {
	return &RepoAuthentication{
		db: db,
	}
}

type RepoAuthentication struct {
	db model.IDatabase
}

func (authenticationRepo *RepoAuthentication) CreateAccount(ctx context.Context,
	login,
	email,
	password string,
) (int, error) {
	args := pgx.NamedArgs{
		"login":    login,
		"email":    email,
		"password": password,
		"role":     model.RoleUser,
	}
	accountID, err := authenticationRepo.db.Insert(ctx, CreateAccount, args)
	if err != nil {
		return 0, err
	}
	return accountID, nil
}

func (authenticationRepo *RepoAuthentication) SetTwoFaKey(
	ctx context.Context,
	accountID int,
	twoFaKey string,
) error {
	args := pgx.NamedArgs{
		"account_id": accountID,
		"two_fa_key": twoFaKey,
	}
	err := authenticationRepo.db.Update(ctx, SetTwoFaKey, args)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (authenticationRepo *RepoAuthentication) AccountByLogin(ctx context.Context,
	login string,
) ([]*model.Account, error) {
	args := pgx.NamedArgs{
		"login": login,
	}
	rows, err := authenticationRepo.db.Select(ctx, AccountByLogin, args)
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

func (authenticationRepo *RepoAuthentication) AccountByID(ctx context.Context,
	accountID int,
) ([]*model.Account, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	rows, err := authenticationRepo.db.Select(ctx, AccountByID, args)
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

func (authenticationRepo *RepoAuthentication) UpgradeToAdmin(ctx context.Context,
	accountID int,
) error {
	args := pgx.NamedArgs{
		"account_id": accountID,
		"role":       model.RoleAdmin,
	}
	err := authenticationRepo.db.Update(ctx, UpgradeToAdmin, args)
	if err != nil {
		return err
	}
	return nil
}

func (authenticationRepo *RepoAuthentication) UpgradeToSupport(ctx context.Context,
	accountID int,
) error {
	args := pgx.NamedArgs{
		"account_id": accountID,
		"role":       model.RoleSupport,
	}
	err := authenticationRepo.db.Update(ctx, UpgradeToSupport, args)
	if err != nil {
		return err
	}
	return nil
}

func (authenticationRepo *RepoAuthentication) UpdatePassword(ctx context.Context,
	accountID int,
	newPassword string,
) error {
	args := pgx.NamedArgs{
		"account_id":              accountID,
		"new_password":            newPassword,
		"last_change_password_at": time.Now(),
	}
	err := authenticationRepo.db.Update(ctx, UpdatePassword, args)
	if err != nil {
		return err
	}
	return nil
}

func (authenticationRepo *RepoAuthentication) DeleteTwoFaKey(ctx context.Context,
	accountID int,
) error {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	err := authenticationRepo.db.Update(ctx, DeleteTwoFaKey, args)
	if err != nil {
		return err
	}
	return nil
}
