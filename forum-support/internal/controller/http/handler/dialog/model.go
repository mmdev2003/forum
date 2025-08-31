package dialog

type CreateDialogRequest struct {
	SupportRequestID int `json:"supportRequestID"`
	UserAccountID    int `json:"userAccountID"`
}
