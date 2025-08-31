package authentication

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"forum-authentication/internal/model"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
	"image/png"
)

func New(
	authenticationRepo model.IAuthenticationRepo,
	authorizationClient model.IAuthorizationClient,
	userClient model.IUserClient,
	adminClient model.IAdminClient,
	threadClient model.IThreadClient,
	passwordSecretKey string,
) *ServiceAuthentication {
	return &ServiceAuthentication{
		authenticationRepo:  authenticationRepo,
		authorizationClient: authorizationClient,
		userClient:          userClient,
		adminClient:         adminClient,
		threadClient:        threadClient,
		passwordSecretKey:   passwordSecretKey,
	}
}

type ServiceAuthentication struct {
	authenticationRepo  model.IAuthenticationRepo
	authorizationClient model.IAuthorizationClient
	userClient          model.IUserClient
	adminClient         model.IAdminClient
	threadClient        model.IThreadClient
	passwordSecretKey   string
}

func (s *ServiceAuthentication) Register(ctx context.Context,
	login,
	email,
	password string,
) (*model.AuthorizationData, error) {
	hashedPassword := s.hashPassword(password)

	accountID, err := s.authenticationRepo.CreateAccount(ctx, login, email, hashedPassword)
	if err != nil {
		return nil, err
	}

	jwtTokens, err := s.authorizationClient.Authorization(ctx, accountID)
	if err != nil {
		return nil, err
	}

	err = s.threadClient.CreateAccountStatistic(ctx, accountID)
	if err != nil {
		return nil, err
	}

	err = s.userClient.CreateUser(ctx, accountID, login)
	if err != nil {
		return nil, err
	}

	return &model.AuthorizationData{
		AccountID:       accountID,
		Role:            "user",
		AccessToken:     jwtTokens.AccessToken,
		RefreshToken:    jwtTokens.RefreshToken,
		IsTwoFaVerified: true,
	}, nil
}

func (s *ServiceAuthentication) Login(ctx context.Context,
	login,
	password,
	twoFaCode string,
) (*model.AuthorizationData, error) {
	account, err := s.authenticationRepo.AccountByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if len(account) == 0 {
		return nil, model.ErrAccountNotFound
	}

	err = s.verifyPassword(account[0].Password, password)
	if err != nil {
		return nil, model.ErrInvalidPassword
	}

	jwtTokens, err := s.authorizationClient.Authorization(ctx,
		account[0].ID,
	)
	if err != nil {
		return nil, err
	}

	var isTwoFaVerified bool
	if account[0].TwoFaKey == "" {
		isTwoFaVerified = true
	} else {
		isTwoFaVerified = s.verifyTwoFa(
			account[0].TwoFaKey,
			twoFaCode,
		)
	}

	return &model.AuthorizationData{
		AccountID:            account[0].ID,
		Role:                 account[0].Role,
		LastChangePasswordAt: account[0].LastChangePasswordAt,
		AccessToken:          jwtTokens.AccessToken,
		RefreshToken:         jwtTokens.RefreshToken,
		IsTwoFaVerified:      isTwoFaVerified,
	}, nil

}

func (s *ServiceAuthentication) GenerateTwoFa(
	accountID int,
) (string, *bytes.Buffer, error) {
	twoFaKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "forum",
		AccountName: fmt.Sprintf("account-%d", accountID),
	})
	if err != nil {
		return "", nil, err
	}

	provisioningUri := twoFaKey.URL()

	qrCode, err := qrcode.New(provisioningUri, qrcode.Medium)
	if err != nil {
		return "", nil, err
	}

	var qrBuffer bytes.Buffer
	err = png.Encode(&qrBuffer, qrCode.Image(256))
	if err != nil {
		return "", nil, err
	}

	return twoFaKey.Secret(), &qrBuffer, nil
}

func (s *ServiceAuthentication) SetTwoFaKey(ctx context.Context,
	accountID int,
	twoFaKey string,
	twoFaCode string,
) error {
	account, err := s.authenticationRepo.AccountByID(ctx, accountID)
	if err != nil {
		return err
	}
	if len(account) == 0 {
		return model.ErrAccountNotFound
	}
	if account[0].TwoFaKey != "" {
		return model.ErrTwoFaAlreadyEnabled
	}

	verified := s.verifyTwoFa(twoFaKey, twoFaCode)
	if !verified {
		return model.ErrTwoFaCodeInvalid
	}
	err = s.authenticationRepo.SetTwoFaKey(ctx, accountID, twoFaKey)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceAuthentication) DeleteTwoFaKey(ctx context.Context,
	accountID int,
	twoFaCode string,
) error {
	account, err := s.authenticationRepo.AccountByID(ctx, accountID)
	if err != nil {
		return err
	}
	if len(account) == 0 {
		return model.ErrAccountNotFound
	}

	verified, err := s.VerifyTwoFa(ctx, accountID, twoFaCode)
	if err != nil {
		return err
	}
	if !verified {
		return errors.New("two fa not verified")
	}

	err = s.authenticationRepo.DeleteTwoFaKey(ctx, accountID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceAuthentication) VerifyTwoFa(ctx context.Context,
	accountID int,
	twoFaCode string,
) (bool, error) {
	account, err := s.authenticationRepo.AccountByID(ctx, accountID)
	if err != nil {
		return false, err
	}
	if len(account) == 0 {
		return false, model.ErrAccountNotFound
	}
	if account[0].TwoFaKey == "" {
		return false, model.ErrTwoFaNotEnabled
	}
	isTwoFaVerified := s.verifyTwoFa(account[0].TwoFaKey, twoFaCode)
	return isTwoFaVerified, nil
}

func (s *ServiceAuthentication) UpgradeToAdmin(ctx context.Context,
	accountID int,
) error {
	err := s.authenticationRepo.UpgradeToAdmin(ctx, accountID)
	if err != nil {
		return err
	}

	err = s.adminClient.CreateAdmin(ctx, accountID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceAuthentication) UpgradeToSupport(ctx context.Context,
	accountID int,
) error {
	err := s.authenticationRepo.UpgradeToSupport(ctx, accountID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceAuthentication) RecoveryPassword(ctx context.Context,
	accountID int,
	twoFaCode,
	newPassword string,
) error {
	verified, err := s.VerifyTwoFa(ctx, accountID, twoFaCode)
	if err != nil {
		return err
	}
	if !verified {
		return errors.New("two fa not verified")
	}

	hashedPassword := s.hashPassword(newPassword)
	err = s.authenticationRepo.UpdatePassword(ctx, accountID, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceAuthentication) ChangePassword(ctx context.Context,
	accountID int,
	oldPassword,
	newPassword string,
) error {
	account, err := s.authenticationRepo.AccountByID(ctx, accountID)
	if err != nil {
		return err
	}
	if len(account) == 0 {
		return model.ErrAccountNotFound
	}

	err = s.verifyPassword(account[0].Password, oldPassword)
	if err != nil {
		return model.ErrInvalidPassword
	}

	hashedPassword := s.hashPassword(newPassword)
	err = s.authenticationRepo.UpdatePassword(ctx, accountID, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceAuthentication) verifyTwoFa(twoFaKey string, twoFaCode string) bool {
	return totp.Validate(twoFaCode, twoFaKey)
}
func (s *ServiceAuthentication) hashPassword(password string) string {
	pepperedPassword := s.passwordSecretKey + password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pepperedPassword), bcrypt.DefaultCost)
	return string(hashedPassword)
}
func (s *ServiceAuthentication) verifyPassword(hashedPassword, password string) error {
	pepperedPassword := s.passwordSecretKey + password
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(pepperedPassword))
}
