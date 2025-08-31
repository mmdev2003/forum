package authentication

import (
	"context"
	"forum-authentication/internal/model"
	"forum-authentication/pkg/api/forum-authorization"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func Register(authenticationService model.IAuthenticationService, domain string) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body RegisterBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		authorizationData, err := authenticationService.Register(
			ctx,
			body.Login,
			body.Email,
			body.Password,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		accessCookie := &http.Cookie{
			Name:     "Access-Token",
			Value:    authorizationData.AccessToken,
			Expires:  time.Now().Add(time.Minute * 15),
			HttpOnly: true,
			Path:     "/",
			Domain:   domain,
		}
		refreshCookie := &http.Cookie{
			Name:     "Refresh-Token",
			Value:    authorizationData.RefreshToken,
			Expires:  time.Now().Add(time.Hour * 1),
			HttpOnly: true,
			Path:     "/",
			Domain:   domain,
		}

		request.SetCookie(refreshCookie)
		request.SetCookie(accessCookie)
		return request.JSON(http.StatusOK, RegisterResponse{AccountID: authorizationData.AccountID})
	}
}

func Login(authenticationService model.IAuthenticationService, domain string) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body LoginBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		authorizationData, err := authenticationService.Login(
			ctx,
			body.Login,
			body.Password,
			body.TwoFaCode,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if authorizationData.IsTwoFaVerified {
			accessCookie := &http.Cookie{
				Name:     "Access-Token",
				Value:    authorizationData.AccessToken,
				Expires:  time.Now().Add(time.Minute * 15),
				HttpOnly: true,
				Path:     "/",
				Domain:   domain,
			}
			refreshCookie := &http.Cookie{
				Name:     "Refresh-Token",
				Value:    authorizationData.RefreshToken,
				Expires:  time.Now().Add(time.Hour * 1),
				HttpOnly: true,
				Path:     "/",
				Domain:   domain,
			}

			request.SetCookie(refreshCookie)
			request.SetCookie(accessCookie)
		}
		return request.JSON(http.StatusOK, LoginResponse{
			AccountID:            authorizationData.AccountID,
			Role:                 authorizationData.Role,
			IsTwoFaVerified:      authorizationData.IsTwoFaVerified,
			LastChangePasswordAt: authorizationData.LastChangePasswordAt,
		})

	}
}

func GenerateTwoFa(authenticationService model.IAuthenticationService) echo.HandlerFunc {
	return func(request echo.Context) error {
		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		twoFaKey, qrImage, err := authenticationService.GenerateTwoFa(authorizationData.AccountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		request.Response().Header().Set("Two-Fa-Key", twoFaKey)

		return request.Blob(200, "image/png", qrImage.Bytes())
	}
}

func SetTwoFa(authenticationService model.IAuthenticationService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body SetTwoFaBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)
		if authorizationData.AccountID == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, model.ErrUnauthorized)
		}

		err := authenticationService.SetTwoFaKey(ctx,
			authorizationData.AccountID,
			body.TwoFaKey,
			body.TwoFaCode,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, "ok")
	}
}

func DeleteTwoFaKey(authenticationService model.IAuthenticationService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		twoFaCode := request.Param("twoFaCode")

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)
		if authorizationData.AccountID == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, model.ErrUnauthorized)
		}

		err := authenticationService.DeleteTwoFaKey(ctx,
			authorizationData.AccountID,
			twoFaCode,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, "ok")
	}
}

func VerifyTwoFa(authenticationService model.IAuthenticationService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := context.Background()
		twoFaCode := request.Param("twoFaCode")

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)
		if authorizationData.AccountID == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, model.ErrUnauthorized)
		}

		isTwoFaVerified, err := authenticationService.VerifyTwoFa(ctx, authorizationData.AccountID, twoFaCode)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, VerifyTwoFaResponse{isTwoFaVerified})
	}
}

func UpgradeToAdmin(authenticationService model.IAuthenticationService, interServerSecretKey string) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := context.Background()

		var body UpgradeToAdminBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		if body.InterServerSecretKey != interServerSecretKey {
			return echo.NewHTTPError(http.StatusUnauthorized, model.ErrUnauthorized)
		}

		err := authenticationService.UpgradeToAdmin(ctx, body.AccountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "ok")
	}
}

func UpgradeToSupport(authenticationService model.IAuthenticationService, interServerSecretKey string) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := context.Background()

		var body UpgradeToSupportBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		if body.InterServerSecretKey != interServerSecretKey {
			return echo.NewHTTPError(http.StatusUnauthorized, model.ErrUnauthorized)
		}

		err := authenticationService.UpgradeToSupport(ctx, body.AccountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "ok")
	}
}

func RecoveryPassword(authenticationService model.IAuthenticationService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := context.Background()

		var body RecoveryPasswordBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)
		if authorizationData.AccountID == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, model.ErrUnauthorized)
		}

		err := authenticationService.RecoveryPassword(ctx, authorizationData.AccountID, body.TwoFaCode, body.NewPassword)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "ok")
	}
}

func ChangePassword(authenticationService model.IAuthenticationService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := context.Background()

		var body ChangePasswordBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)
		if authorizationData.AccountID == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, model.ErrUnauthorized)
		}

		err := authenticationService.ChangePassword(ctx, authorizationData.AccountID, body.OldPassword, body.NewPassword)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "ok")
	}
}
