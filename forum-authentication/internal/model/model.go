package model

import "time"

type Account struct {
	ID int `db:"id"`

	Login    string `db:"login"`
	Email    string `db:"email"`
	Password string `db:"password"`

	Role     string `db:"role"`
	TwoFaKey string `db:"two_fa_key"`

	LastChangePasswordAt time.Time `db:"last_change_password_at"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
}

type AuthorizationData struct {
	AccountID            int       `json:"accountID"`
	Role                 string    `json:"role"`
	IsTwoFaVerified      bool      `json:"isTwoFaVerified"`
	LastChangePasswordAt time.Time `json:"lastChangePasswordAt"`
	AccessToken          string    `json:"accessToken"`
	RefreshToken         string    `json:"refreshToken"`
}
