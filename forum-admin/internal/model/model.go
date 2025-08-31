package model

import "time"

type Admin struct {
	ID        int       `db:"id"`
	AccountID int       `db:"account_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
