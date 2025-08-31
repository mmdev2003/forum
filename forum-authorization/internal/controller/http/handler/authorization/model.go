package authorization

type AuthBody struct {
	AccountID   int    `json:"accountID"`
	Role        string `json:"role"`
	TwoFaStatus bool   `json:"twoFaStatus"`
}

type CheckAuthorizationResponse struct {
	AccountID   int    `json:"accountID"`
	Role        string `json:"role"`
	TwoFaStatus bool   `json:"twoFaStatus"`

	Message string `json:"message"`
	Code    int    `json:"code"`
}
