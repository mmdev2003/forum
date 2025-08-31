package user

import (
	"context"
	"encoding/json"
	"forum-user/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"strconv"
)

func New(
	db model.IDatabase,
	fullTextSearchEngine model.IFullTextSearchEngine,
	weedFS model.IWeedFS,
) *RepoUser {
	return &RepoUser{
		db:                   db,
		fullTextSearchEngine: fullTextSearchEngine,
		weedFS:               weedFS,
	}
}

type RepoUser struct {
	db                   model.IDatabase
	fullTextSearchEngine model.IFullTextSearchEngine
	weedFS               model.IWeedFS
}

func (r *RepoUser) CreateUser(ctx context.Context,
	accountID int,
	login string,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
		"login":      login,
	}
	userID, err := r.db.Insert(ctx, CreateUser, args)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *RepoUser) CreateUserBan(ctx context.Context,
	fromAccountID,
	toAccountID int,
) error {
	args := pgx.NamedArgs{
		"from_account_id": fromAccountID,
		"to_account_id":   toAccountID,
	}
	_, err := r.db.Insert(ctx, CreateUserBan, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoUser) CreateWarningFromAdmin(ctx context.Context,
	adminAccountID,
	toAccountID int,
	warningType,
	warningText,
	adminLogin string,
) error {
	args := pgx.NamedArgs{
		"admin_account_id": adminAccountID,
		"to_account_id":    toAccountID,
		"warning_type":     warningType,
		"warning_text":     warningText,
		"admin_login":      adminLogin,
	}
	_, err := r.db.Insert(ctx, CreateWarningFromAdmin, args)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoUser) CreateUserSearchIndex(ctx context.Context,
	accountID int,
	login string,
) error {
	document := model.UserSearch{
		ID:    accountID,
		Login: login,
	}
	err := r.fullTextSearchEngine.AddDocuments(model.UserFullTextSearchIndex, []any{document})
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoUser) UploadAvatar(ctx context.Context,
	accountID int,
	avatar []byte,
) (string, error) {
	accountIDStr := strconv.Itoa(accountID)
	fileID, err := r.weedFS.Upload(avatar, "avatar-"+accountIDStr, int64(len(avatar)), "")
	if err != nil {
		return "", err
	}
	return fileID, nil
}

func (r *RepoUser) UploadHeader(ctx context.Context,
	accountID int,
	header []byte,
) (string, error) {
	accountIDStr := strconv.Itoa(accountID)
	fileID, err := r.weedFS.Upload(header, "header-"+accountIDStr, int64(len(header)), "")
	if err != nil {
		return "", err
	}
	return fileID, nil
}

func (r *RepoUser) UpdateAvatarUrl(ctx context.Context,
	accountID int,
	avatarUrl string,
) error {
	args := pgx.NamedArgs{
		"account_id": accountID,
		"avatar_url": avatarUrl,
	}
	err := r.db.Update(ctx, UpdateAvatarUrl, args)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoUser) UpdateHeaderUrl(ctx context.Context,
	accountID int,
	headerUrl string,
) error {
	args := pgx.NamedArgs{
		"account_id": accountID,
		"header_url": headerUrl,
	}
	err := r.db.Update(ctx, UpdateHeaderUrl, args)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoUser) DownloadAvatar(ctx context.Context,
	fileID string,
) ([]byte, error) {
	avatar, err := r.weedFS.Download(fileID)
	if err != nil {
		return nil, err
	}
	return avatar, nil
}

func (r *RepoUser) DownloadHeader(ctx context.Context,
	fileID string,
) ([]byte, error) {
	header, err := r.weedFS.Download(fileID)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (r *RepoUser) UserByAccountID(ctx context.Context,
	accountID int,
) ([]*model.User, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	rows, err := r.db.Select(ctx, GetUserByAccountID, args)
	if err != nil {
		return nil, err
	}

	var users []*model.User
	err = pgxscan.ScanAll(&users, rows)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *RepoUser) UserByLogin(ctx context.Context,
	login string,
) ([]*model.User, error) {
	args := pgx.NamedArgs{
		"login": login,
	}
	rows, err := r.db.Select(ctx, GetUserByLogin, args)
	if err != nil {
		return nil, err
	}

	var user []*model.User
	err = pgxscan.ScanAll(&user, rows)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *RepoUser) UsersByLoginSearch(ctx context.Context,
	login string,
) ([]*model.UserSearch, error) {
	result, err := r.fullTextSearchEngine.SimpleSearch(model.UserFullTextSearchIndex, login)
	if err != nil {
		return nil, err
	}

	var users []*model.UserSearch
	err = json.Unmarshal(result, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *RepoUser) BanByAccountID(ctx context.Context,
	toAccountID int,
) ([]*model.UserBan, error) {
	args := pgx.NamedArgs{
		"to_account_id": toAccountID,
	}
	rows, err := r.db.Select(ctx, BanByAccountID, args)
	if err != nil {
		return nil, err
	}

	var userBans []*model.UserBan
	err = pgxscan.ScanAll(&userBans, rows)
	if err != nil {
		return nil, err
	}

	return userBans, nil
}

func (r *RepoUser) AllWarningFromAdmin(ctx context.Context,
	toAccountID int,
) ([]*model.WarningFromAdmin, error) {
	args := pgx.NamedArgs{
		"to_account_id": toAccountID,
	}
	rows, err := r.db.Select(ctx, AllWarningFromAdmin, args)
	if err != nil {
		return nil, err
	}

	var userWarnings []*model.WarningFromAdmin
	err = pgxscan.ScanAll(&userWarnings, rows)
	if err != nil {
		return nil, err
	}

	return userWarnings, nil
}

func (r *RepoUser) DeleteUserBan(ctx context.Context,
	fromAccountID,
	toAccountID int,
) error {
	args := pgx.NamedArgs{
		"from_account_id": fromAccountID,
		"to_account_id":   toAccountID,
	}
	err := r.db.Delete(ctx, DeleteUserBan, args)
	if err != nil {
		return err
	}

	return nil
}
