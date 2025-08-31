package model

import "time"

type SupportRequest struct {
	ID          int       `db:"id" json:"id"`
	AccountID   int       `db:"account_id" json:"accountID"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Status      *string   `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

type Dialog struct {
	ID               int        `db:"id" json:"id"`
	SupportRequestID int        `db:"support_request_id" json:"supportRequestID"`
	UserAccountID    int        `db:"user_account_id" json:"userAccountID"`
	LastMessageAt    *time.Time `db:"last_message_at" json:"lastMessageAt"`
	CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updatedAt"`
}

type Message struct {
	ID            int       `db:"id" json:"id"`
	DialogID      int       `db:"dialog_id" json:"dialogID"`
	FromAccountID int       `db:"from_account_id" json:"fromAccountID"`
	ToAccountID   int       `db:"to_account_id" json:"toAccountID"`
	MessageText   string    `db:"message_text" json:"messageText"`
	IsRead        bool      `db:"is_read" json:"isRead"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
}

type DialogWsMessage struct {
	DialogID      int    `json:"dialogID"`
	FromAccountID int    `json:"fromAccountID"`
	ToAccountID   int    `json:"toAccountID"`
	Text          string `json:"text"`
}
