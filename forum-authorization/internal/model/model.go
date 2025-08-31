package model

import "time"

type Account struct {
	ID int `db:"id"`

	AccountID    int    `db:"account_id"`
	Role         string `db:"role"`
	TwoFaKey     string `db:"two_fa_key"`
	RefreshToken string `db:"refresh_token"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type JWTTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenPayload struct {
	AccountID   int
	Role        string
	TwoFaStatus bool
	Exp         int64
}
