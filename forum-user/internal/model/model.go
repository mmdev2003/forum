package model

import "time"

type User struct {
	ID        int       `db:"id"`
	AccountID int       `db:"account_id"`
	Login     string    `db:"login"`
	HeaderUrl string    `db:"header_url"`
	AvatarUrl string    `db:"avatar_url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserBan struct {
	ID            int       `db:"id"`
	FromAccountID int       `db:"from_account_id"`
	ToAccountID   int       `db:"to_account_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type WarningFromAdmin struct {
	ID             int       `db:"id"`
	AdminAccountID int       `db:"admin_account_id"`
	AdminLogin     string    `db:"admin_login"`
	ToAccountID    int       `db:"to_account_id"`
	WarningType    string    `db:"warning_type"`
	WarningText    string    `db:"warning_text"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

type UserSearch struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}
