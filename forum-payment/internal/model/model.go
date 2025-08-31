package model

import (
	"time"
)

type Payment struct {
	ID        int `db:"id"`
	AccountID int `db:"account_id"`

	ProductType ProductType   `db:"product_type"`
	Address     string        `db:"address"`
	Currency    Currency      `db:"currency"`
	Amount      string        `db:"amount"`
	Status      PaymentStatus `db:"status"`
	TxID        string        `db:"tx_id"`
	IsPaid      bool          `db:"is_paid"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
