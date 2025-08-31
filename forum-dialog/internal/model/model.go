package model

import "time"

type Dialog struct {
	ID                  int        `db:"id"`
	Account1ID          int        `db:"account1_id"`
	Account2ID          int        `db:"account2_id"`
	IsStarredByAccount1 bool       `db:"is_starred_by_account1"`
	IsStarredByAccount2 bool       `db:"is_starred_by_account2"`
	LastMessageAt       *time.Time `db:"last_message_at"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
}

type Message struct {
	ID            int       `db:"id"`
	DialogID      int       `db:"dialog_id"`
	FromAccountID int       `db:"from_account_id"`
	ToAccountID   int       `db:"to_account_id"`
	MessageText   string    `db:"message_text"`
	IsRead        bool      `db:"is_read"`
	FileID        int       `db:"file_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type File struct {
	ID        int       `db:"id"`
	MessageID int       `db:"message_id"`
	URL       string    `db:"url"`
	Name      string    `db:"name"`
	Size      int       `db:"size"`
	Extension string    `db:"extension"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type DialogWsMessage struct {
	DialogID      int      `json:"dialogID"`
	FromAccountID int      `json:"fromAccountID"`
	ToAccountID   int      `json:"toAccountID"`
	Text          string   `json:"text"`
	FilesURLs     []string `json:"filesURLs"`
}
