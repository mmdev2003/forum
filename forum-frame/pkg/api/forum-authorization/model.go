package forum_authorization

type JWTTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthorizationData struct {
	AccountID   int    `json:"accountID"`
	Role        string `json:"role"`
	TwoFaStatus bool   `json:"twoFaStatus"`
	Message     string `json:"message"`
	Code        int    `json:"code"`
}
