package dialog

import (
	"context"
	"forum-support/internal/model"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type RepoDialog struct {
	db model.IDatabase
}

func New(db model.IDatabase) *RepoDialog {
	return &RepoDialog{db: db}
}

func (r RepoDialog) insert(ctx context.Context, query string, args any) (int, error) {
	id, err := r.db.Insert(ctx, query, args)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r RepoDialog) CreateDialog(ctx context.Context, requestAccountID, userAccountID int) (int, error) {
	args := pgx.NamedArgs{
		"support_request_id": requestAccountID,
		"user_account_id":    userAccountID,
	}
	return r.insert(ctx, CreateDialog, args)
}

func (r RepoDialog) AddMessageToDialog(ctx context.Context, dialogID, fromAccountID, toAccountID int, text string) (int, error) {
	args := pgx.NamedArgs{
		"dialog_id":       dialogID,
		"from_account_id": fromAccountID,
		"to_account_id":   toAccountID,
		"message_text":    text,
	}
	return r.insert(ctx, AddMessageToDialog, args)
}

func (r RepoDialog) MarkMessagesAsRead(ctx context.Context, dialogID int, requesterID int, as_support bool) error {
	args := pgx.NamedArgs{
		"dialog_id":    dialogID,
		"requester_id": requesterID,
	}
	if as_support {
		return r.db.Update(ctx, MarkMessagesFromRequesterAsRead, args)
	} else {
		return r.db.Update(ctx, MarkMessagesFromSupportAsRead, args)
	}
}

func (r RepoDialog) AllDialogs(ctx context.Context) ([]*model.Dialog, error) {
	rows, err := r.db.Select(ctx, GetAllDialogs)
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

func (r RepoDialog) GetDialogByID(ctx context.Context, dialogID int) (*model.Dialog, error) {
	args := pgx.NamedArgs{
		"dialog_id": dialogID,
	}

	rows, err := r.db.Select(ctx, GetDialogByID, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, pgx.ErrNoRows
	}

	var dialog model.Dialog
	err = rows.Scan(
		&dialog.ID,
		&dialog.SupportRequestID,
		&dialog.UserAccountID,
		&dialog.LastMessageAt,
		&dialog.CreatedAt,
		&dialog.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &dialog, nil
}

func (r RepoDialog) DialogsByAccountID(ctx context.Context, accountID int) ([]*model.Dialog, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}

	rows, err := r.db.Select(ctx, GetDialogsByAccountID, args)
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

func (r RepoDialog) MessagesByDialogID(ctx context.Context, dialogID int) ([]*model.Message, error) {
	args := pgx.NamedArgs{
		"dialog_id": dialogID,
	}
	rows, err := r.db.Select(ctx, MessagesByDialogId, args)
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
