package handler

import (
	"forum-user/internal/model"
	"time"
)

type CreateUserBody struct {
	AccountID int    `json:"accountID" validate:"required"`
	Login     string `json:"login"`
}

type BanUserBody struct {
	ToAccountID int `json:"to_account_id" validate:"required"`
}

type NewWarningFromAdminBody struct {
	ToAccountID int    `json:"to_account_id" validate:"required"`
	WarningType string `json:"warning_type" validate:"required"`
	AdminLogin  string `json:"admin_login" validate:"required"`
}

type GetUserByAccountIDBody struct {
	AccountID int `json:"accountID" validate:"required"`
}

type GetUserByAccountIDResponse struct {
	AccountID int `json:"account_id"`

	Login     string `json:"login"`
	HeaderUrl string `json:"header_url"`
	AvatarUrl string `json:"avatar_url"`

	CreatedAt time.Time `json:"created_at"`
}

type UsersByLoginBody struct {
	Login string `json:"login" validate:"required"`
}

type UsersByLoginResponse struct {
	Users []*model.UserSearch `json:"users"`
}

type GetUserByLoginResponse struct {
	User []*model.User `json:"user"`
}
type BanByAccountIDResponse struct {
	UserBans []*model.UserBan `json:"user_bans"`
}

type AllWarningFromAdminResponse struct {
	UserWarnings []*model.WarningFromAdmin `json:"user_warnings"`
}
