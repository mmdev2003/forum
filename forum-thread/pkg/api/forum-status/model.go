package forum_status

import "time"

type AssignStatusToAccountBody struct {
	AccountID            int    `json:"accountID"`
	StatusID             int    `json:"statusID"`
	InterServerSecretKey string `json:"interServerSecretKey"`
}

type Status struct {
	ID            int    `db:"id"`
	StatusID      int    `db:"status_id"`
	AccountID     int    `db:"account_id"`
	PaymentID     int    `db:"payment_id"`
	PaymentStatus string `db:"payment_status"`

	ExpirationAt time.Time `db:"expiration_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type StatusByAccountIDResponse struct {
	Statuses []*Status `json:"statuses"`
}
