package authorization

import (
	"errors"
	"forum-authorization/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func Authorization(authorizationService model.IAuthorizationService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body AuthBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		JWTTokens, err := authorizationService.CreateTokens(
			ctx,
			body.AccountID,
			body.Role,
			body.TwoFaStatus,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(200, JWTTokens)
	}
}

func CheckAuthorization(authorizationService model.IAuthorizationService) echo.HandlerFunc {
	return func(request echo.Context) error {
		accessToken, err := request.Cookie("Access-Token")

		var response CheckAuthorizationResponse
		if err != nil {
			response.AccountID = 0
			response.TwoFaStatus = false
			response.Role = "Guest"
			response.Message = "Access-Token non set"
			response.Code = model.CodeErrAccessTokenNonSet
			return echo.NewHTTPError(http.StatusUnauthorized, response)
		}
		tokenPayload, err := authorizationService.CheckToken(accessToken.Value)
		if err != nil {
			response.AccountID = -1
			response.TwoFaStatus = false
			response.Role = "Guest"
			switch {
			case errors.Is(err, model.ErrTokenExpired):
				response.Message = "Access-Token expired"
				response.Code = model.CodeErrAccessTokenExpired

			case errors.Is(err, model.ErrTokenInvalid):
				response.Message = "Access-Token invalid"
				response.Code = model.CodeErrAccessTokenInvalid

			}
			return echo.NewHTTPError(http.StatusUnauthorized, response)
		}
		response.AccountID = tokenPayload.AccountID
		response.TwoFaStatus = tokenPayload.TwoFaStatus
		response.Role = tokenPayload.Role
		response.Message = "Access-Token verified"
		response.Code = 200

		return request.JSON(http.StatusOK, response)
	}
}

func RefreshTokens(authorizationService model.IAuthorizationService, domain string) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		refreshToken, err := request.Cookie("Refresh-Token")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		JWTTokens, err := authorizationService.RefreshTokens(ctx, refreshToken.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		accessCookie := &http.Cookie{
			Name:     "Access-Token",
			Value:    JWTTokens.AccessToken,
			Expires:  time.Now().Add(time.Minute * 15),
			HttpOnly: true,
			Path:     "/",
			Domain:   domain,
		}
		refreshCookie := &http.Cookie{
			Name:     "Refresh-Token",
			Value:    JWTTokens.RefreshToken,
			Expires:  time.Now().Add(time.Hour * 1),
			HttpOnly: true,
			Path:     "/",
			Domain:   domain,
		}

		request.SetCookie(refreshCookie)
		request.SetCookie(accessCookie)

		return request.String(http.StatusOK, "ok")
	}
}
