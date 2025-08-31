package forum_user

type CreateUserBody struct {
	AccountID int    `json:"accountID"`
	Login     string `json:"login"`
}
