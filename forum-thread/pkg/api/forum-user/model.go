package forum_user

import "time"

type User struct {
	ID        int `json:"id"`
	AccountID int `json:"account_id"`

	Login     string `json:"login"`
	HeaderUrl string `json:"header_url"`
	AvatarUrl string `json:"avatar_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetUserByLoginResponse struct {
	User []*User `db:"user"`
}
