package dialog

import (
	"context"
	"encoding/json"
	"forum-support/internal/model"
)

type ServiceDialog struct {
	dialogRepo    model.IDialogRepo
	wsConnManager model.IWsConnManager
}

func New(dialogRepo model.IDialogRepo, wsConnManager model.IWsConnManager) *ServiceDialog {
	return &ServiceDialog{dialogRepo: dialogRepo, wsConnManager: wsConnManager}
}

func (s ServiceDialog) CreateDialog(ctx context.Context, supportRequestID, userAccountID int) (int, error) {
	return s.dialogRepo.CreateDialog(ctx, supportRequestID, userAccountID)
}

func (s ServiceDialog) AddMessageToDialog(ctx context.Context, dialogID, fromAccountID, toAccountID int, text string) (int, error) {
	messageId, err := s.dialogRepo.AddMessageToDialog(ctx, dialogID, fromAccountID, toAccountID, text)
	if err != nil {
		return messageId, err
	}

	wsBody, err := json.Marshal(model.DialogWsMessage{
		DialogID:      dialogID,
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Text:          text,
	})

	if err != nil {
		return messageId, err
	}

	err = s.wsConnManager.SendMsg(ctx, toAccountID, wsBody)

	return messageId, err
}

func (s ServiceDialog) MarkMessagesAsRead(ctx context.Context, dialogID int, accountID int, role string) error {
	dialog, err := s.dialogRepo.GetDialogByID(ctx, dialogID)
	if err != nil {
		return err
	}

	if role != model.RoleSupport && dialog.UserAccountID != accountID {
		return model.ErrActionNotAllowed
	}

	return s.dialogRepo.MarkMessagesAsRead(ctx, dialogID, dialog.UserAccountID, role == model.RoleSupport)
}

func (s ServiceDialog) AllDialogs(ctx context.Context, accountID int, role string) ([]*model.Dialog, error) {
	if role == model.RoleSupport {
		return s.dialogRepo.AllDialogs(ctx)
	} else {
		return s.dialogRepo.DialogsByAccountID(ctx, accountID)
	}
}

func (s ServiceDialog) MessagesByDialogID(ctx context.Context, dialogID int, accountID int, role string) ([]*model.Message, error) {
	if role != model.RoleSupport {
		dialog, err := s.dialogRepo.GetDialogByID(ctx, dialogID)
		if err != nil {
			return nil, err
		}

		if dialog.UserAccountID != accountID {
			return nil, model.ErrActionNotAllowed
		}
	}

	return s.dialogRepo.MessagesByDialogID(ctx, dialogID)
}
