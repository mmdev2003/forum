package forum_admin

type CreateAdminBody struct {
	AccountID            int    `json:"accountID"`
	InterServerSecretKey string `json:"interServerSecretKey"`
}
