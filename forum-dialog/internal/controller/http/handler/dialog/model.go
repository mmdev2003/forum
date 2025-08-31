package dialog

import "forum-dialog/internal/model"

type CreateDialogBody struct {
	Account2ID int `json:"account2ID"`
}

type CreateDialogResponse struct {
	DialogID int `json:"dialogID"`
}

type UploadFileResponse struct {
	FileURL string `json:"fileURL"`
}

type DialogsByAccountIDResponse struct {
	Dialogs []*model.Dialog `json:"dialogs"`
}

type MessagesByDialogIDResponse struct {
	Messages []*model.Message `json:"messages"`
	Files    []*model.File    `json:"files"`
}
