package dialog

import (
	"context"
	"encoding/json"
	"forum-dialog/internal/model"
	"path/filepath"
	"strings"
)

func New(
	dialogRepo model.IDialogRepo,
	wsConnManager model.IWsConnManager,
) *ServiceDialog {
	return &ServiceDialog{
		dialogRepo:    dialogRepo,
		wsConnManager: wsConnManager,
	}
}

type ServiceDialog struct {
	dialogRepo    model.IDialogRepo
	wsConnManager model.IWsConnManager
}

func (s *ServiceDialog) CreateDialog(ctx context.Context,
	account1ID,
	account2ID int,
) (int, error) {
	dialogID, err := s.dialogRepo.CreateDialog(ctx, account1ID, account2ID)
	if err != nil {
		return 0, err
	}
	return dialogID, nil
}

func (s *ServiceDialog) AddMessageToDialog(ctx context.Context,
	dialogID,
	fromAccountID,
	toAccountID int,
	text string,
	filesURLs []string,
) (int, error) {
	ctx, err := s.dialogRepo.CtxWithTx(ctx)
	if err != nil {
		return 0, err
	}

	messageID, err := s.dialogRepo.CreateMessage(ctx, dialogID, fromAccountID, toAccountID, text)
	if err != nil {
		s.dialogRepo.RollbackTx(ctx)
		return 0, err
	}

	for _, fileURL := range filesURLs {
		err = s.dialogRepo.AddFileToMessage(ctx, messageID, fileURL)
		if err != nil {
			return 0, err
		}
	}

	err = s.dialogRepo.UpdateLastMessageAt(ctx, dialogID)
	if err != nil {
		s.dialogRepo.RollbackTx(ctx)
		return 0, err
	}

	err = s.dialogRepo.CommitTx(ctx)
	if err != nil {
		return 0, err
	}

	wsBody, err := json.Marshal(model.DialogWsMessage{
		DialogID:      dialogID,
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Text:          text,
		FilesURLs:     filesURLs,
	})

	err = s.wsConnManager.SendMsg(ctx, toAccountID, wsBody)
	return messageID, nil
}

func (s *ServiceDialog) DeleteDialog(ctx context.Context,
	dialogID int,
) error {
	err := s.dialogRepo.DeleteDialog(ctx, dialogID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceDialog) DialogsByAccountID(ctx context.Context,
	accountID int,
) ([]*model.Dialog, error) {
	dialogs, err := s.dialogRepo.DialogsByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return dialogs, nil
}

func (s *ServiceDialog) MessagesByDialogID(ctx context.Context,
	dialogID int,
) ([]*model.Message, []*model.File, error) {
	messages, err := s.dialogRepo.MessagesByDialogID(ctx, dialogID)
	if err != nil {
		return nil, nil, err
	}

	var files []*model.File
	for _, message := range messages {
		messageFiles, err := s.dialogRepo.FilesByMessageID(ctx, message.ID)
		if err != nil {
			return messages, nil, err
		}
		for _, file := range messageFiles {
			files = append(files, file)
		}
	}
	return messages, files, nil
}

func (s *ServiceDialog) MarkDialogAsStarred(ctx context.Context,
	accountID,
	dialogID int,
) error {
	err := s.dialogRepo.MarkDialogAsStarred(ctx, accountID, dialogID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceDialog) MarkMessagesAsRead(ctx context.Context,
	dialogID int,
) error {
	err := s.dialogRepo.MarkMessagesAsRead(ctx, dialogID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceDialog) DownloadFile(ctx context.Context,
	fileURL string,
) ([]byte, error) {
	file, err := s.dialogRepo.DownloadFile(ctx, fileURL)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *ServiceDialog) UploadFile(ctx context.Context,
	file []byte,
	fullName string,
) (string, error) {
	fileURL, err := s.dialogRepo.UploadFile(ctx, file, fullName)
	if err != nil {
		return "", err
	}

	extension := filepath.Ext(fullName)
	name := strings.TrimSuffix(fullName, extension)
	err = s.dialogRepo.CreateFile(ctx, len(file), fileURL, name, extension)
	if err != nil {
		return "", err
	}

	return fileURL, nil
}
