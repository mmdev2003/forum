package model

import "time"

type Frame struct {
	ID            int           `db:"id"`
	FrameID       int           `db:"frame_id"`
	AccountID     int           `db:"account_id"`
	PaymentID     int           `db:"payment_id"`
	PaymentStatus PaymentStatus `db:"payment_status"`

	ExpirationAt time.Time `db:"expiration_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type FrameData struct {
	ID           int     `db:"id"`
	Name         string  `db:"name"`
	MonthlyPrice float32 `db:"monthly_price"`
	ForeverPrice float32 `db:"forever_price"`
	FileID       string  `db:"file_id"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CurrentFrame struct {
	ID        int       `db:"id"`
	DbFrameID int       `db:"db_frame_id"`
	AccountID int       `db:"account_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
