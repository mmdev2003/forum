package admin

import "forum-admin/internal/model"

type CreateAdminBody struct {
	AccountID            int    `json:"accountID"`
	InterServerSecretKey string `json:"interServerSecretKey"`
}

type AllAdminResponse struct {
	Admins []*model.Admin `json:"admins"`
}
