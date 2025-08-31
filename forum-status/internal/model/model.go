package model

import "time"

type Status struct {
	ID            int           `db:"id"`
	StatusID      int           `db:"status_id"`
	AccountID     int           `db:"account_id"`
	PaymentID     int           `db:"payment_id"`
	PaymentStatus PaymentStatus `db:"payment_status"`

	ExpirationAt time.Time `db:"expiration_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
