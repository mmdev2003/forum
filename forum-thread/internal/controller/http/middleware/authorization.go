package middleware

import (
	"forum-thread/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AuthorizationMiddleware(authorizationClient model.IAuthorizationClient, pathForSkipAuthorization []string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(request echo.Context) error {
			ctx := request.Request().Context()

			for _, path := range pathForSkipAuthorization {
				if request.Request().URL.Path == path {
					return next(request)
				}
			}

			authorizationData, err := authorizationClient.CheckAuthorization(ctx, request)
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
