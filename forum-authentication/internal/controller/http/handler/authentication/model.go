package authentication

import "time"

type RegisterBody struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	AccountID int `json:"accountID"`
}

type LoginBody struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	TwoFaCode string `json:"twoFaCode"`
}

type LoginResponse struct {
	AccountID            int       `json:"accountID"`
	Role                 string    `json:"role"`
	IsTwoFaVerified      bool      `json:"isTwoFaVerified"`
	LastChangePasswordAt time.Time `json:"lastChangePasswordAt"`
}

type SetTwoFaBody struct {
	TwoFaKey  string `json:"twoFaKey"`
	TwoFaCode string `json:"twoFaCode"`
}

type VerifyTwoFaResponse struct {
	IsTwoFaVerified bool `json:"isTwoFaVerified"`
}

type UpgradeToAdminBody struct {
	AccountID            int    `json:"accountID"`
	InterServerSecretKey string `json:"interServerSecretKey"`
}

type UpgradeToSupportBody struct {
	AccountID            int    `json:"accountID"`
	InterServerSecretKey string `json:"interServerSecretKey"`
}

type RecoveryPasswordBody struct {
	TwoFaCode   string `json:"twoFaCode"`
	NewPassword string `json:"newPassword"`
}

type ChangePasswordBody struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
