package dialog

import (
	"context"
)

import (
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

import (
	"forum-dialog/internal/model"
)

func New(
	db model.IDatabase,
	weedFS model.IWeedFS,
) *RepoDialog {
	return &RepoDialog{
		db,
		weedFS,
	}
}

type RepoDialog struct {
	db     model.IDatabase
	weedFS model.IWeedFS
}

func (dialogRepo *RepoDialog) CtxWithTx(ctx context.Context) (context.Context, error) {
	return dialogRepo.db.CtxWithTx(ctx)
}
func (dialogRepo *RepoDialog) CommitTx(ctx context.Context) error {
	return dialogRepo.db.CommitTx(ctx)
}
func (dialogRepo *RepoDialog) RollbackTx(ctx context.Context) {
	dialogRepo.db.RollbackTx(ctx)
}

func (dialogRepo *RepoDialog) CreateDialog(ctx context.Context,
	account1ID,
	account2ID int,
) (int, error) {
	args := pgx.NamedArgs{
		"account1_id": account1ID,
		"account2_id": account2ID,
	}
	dialogID, err := dialogRepo.db.Insert(ctx, CreateDialog, args)
	if err != nil {
		return 0, err
	}

	return dialogID, nil
}

func (dialogRepo *RepoDialog) CreateMessage(ctx context.Context,
	dialogID,
	fromAccountID,
	toAccountID int,
	text string,
) (int, error) {
	args := pgx.NamedArgs{
		"dialog_id":       dialogID,
		"from_account_id": fromAccountID,
		"to_account_id":   toAccountID,
		"message_text":    text,
	}
	messageID, err := dialogRepo.db.Insert(ctx, CreateMessage, args)
	if err != nil {
		return 0, err
	}

	return messageID, nil
}

func (dialogRepo *RepoDialog) UploadFile(ctx context.Context,
	file []byte,
	name string,
) (string, error) {
	fileID, err := dialogRepo.weedFS.Upload(file, name, int64(len(file)), "")
	if err != nil {
		return "", err
	}
	return fileID, nil
}

func (dialogRepo *RepoDialog) CreateFile(ctx context.Context,
	size int,
	fileURL,
	name,
	extension string,
) error {
	args := pgx.NamedArgs{
		"size":      size,
		"url":       fileURL,
		"name":      name,
		"extension": extension,
	}
	_, err := dialogRepo.db.Insert(ctx, CreateFile, args)
	if err != nil {
		return err
	}

	return nil
}

func (dialogRepo *RepoDialog) MarkDialogAsStarred(ctx context.Context,
	accountID,
	dialogID int,
) error {
	args := pgx.NamedArgs{
		"dialog_id":  dialogID,
		"account_id": accountID,
	}
	err := dialogRepo.db.Update(ctx, MarkDialogAsStarred, args)
	if err != nil {
		return err
	}

	return nil
}

func (dialogRepo *RepoDialog) MarkMessagesAsRead(ctx context.Context,
	dialogID int,
) error {
	args := pgx.NamedArgs{
		"dialog_id": dialogID,
	}
	err := dialogRepo.db.Update(ctx, MarkMessagesAsRead, args)
	if err != nil {
		return err
	}

	return nil
}

func (dialogRepo *RepoDialog) UpdateLastMessageAt(ctx context.Context,
	dialogID int,
) error {
	args := pgx.NamedArgs{
		"dialog_id": dialogID,
	}
	err := dialogRepo.db.Update(ctx, UpdateLastMessageAt, args)
	if err != nil {
		return err
	}

	return nil
}

func (dialogRepo *RepoDialog) AddFileToMessage(ctx context.Context,
	messageID int,
	fileURL string,
) error {
	args := pgx.NamedArgs{
		"message_id": messageID,
		"file_url":   fileURL,
	}
	err := dialogRepo.db.Update(ctx, AddFileToMessage, args)
	if err != nil {
		return err
	}

	return nil
}

func (dialogRepo *RepoDialog) DeleteDialog(ctx context.Context,
	dialogID int,
) error {
	args := pgx.NamedArgs{
		"dialog_id": dialogID,
	}
	err := dialogRepo.db.Delete(ctx, DeleteDialog, args)
	if err != nil {
		return err
	}

	return nil
}

func (dialogRepo *RepoDialog) DialogsByAccountID(ctx context.Context,
	accountID int,
) ([]*model.Dialog, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	rows, err := dialogRepo.db.Select(ctx, DialogsByAccountID, args)
	if err != nil {
		return nil, err
	}

	var dialogs []*model.Dialog
	err = pgxscan.ScanAll(&dialogs, rows)
	if err != nil {
		return nil, err
	}
	return dialogs, nil
}

func (dialogRepo *RepoDialog) MessagesByDialogID(ctx context.Context,
	dialogID int,
) ([]*model.Message, error) {
	args := pgx.NamedArgs{
		"dialog_id": dialogID,
	}
	rows, err := dialogRepo.db.Select(ctx, MessagesByDialogId, args)

	if err != nil {
		return nil, err
	}

	var messages []*model.Message
	err = pgxscan.ScanAll(&messages, rows)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
func (dialogRepo *RepoDialog) FilesByMessageID(ctx context.Context,
	messageID int,
) ([]*model.File, error) {
	args := pgx.NamedArgs{
		"message_id": messageID,
	}
	rows, err := dialogRepo.db.Select(ctx, FilesByMessageID, args)
	if err != nil {
		return nil, err
	}

	var files []*model.File
	err = pgxscan.ScanAll(&files, rows)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (dialogRepo *RepoDialog) DownloadFile(ctx context.Context,
	fileURL string,
) ([]byte, error) {
	file, err := dialogRepo.weedFS.Download(fileURL)
	if err != nil {
		return nil, err
	}
	return file, nil
}
