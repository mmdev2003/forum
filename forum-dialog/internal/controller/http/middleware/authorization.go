package middleware

import (
	"forum-dialog/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AuthorizationMiddleware(authClient model.IAuthorizationClient) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(request echo.Context) error {
			ctx := request.Request().Context()

			authorizationData, err := authClient.CheckAuthorization(ctx, request)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, model.ErrCheckAuthorizationFailed)
			}

			switch authorizationData.Code {
			case model.CodeErrAccessTokenExpired:
				return echo.NewHTTPError(http.StatusUnauthorized, model.ErrTokenExpired)
			case model.CodeErrAccessTokenInvalid:
				return echo.NewHTTPError(http.StatusUnauthorized, model.ErrTokenInvalid)
			}

			request.Set(model.AuthorizationDataKey, authorizationData)

			return next(request)
		}
	}
}
