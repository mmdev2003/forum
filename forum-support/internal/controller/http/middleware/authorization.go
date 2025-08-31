package middleware

import (
	"forum-support/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AuthorizationMiddleware(authClient model.IAuthorizationClient) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(request echo.Context) error {
			ctx := request.Request().Context()

			pathForSkipAuthorization := []string{
				"/api/support/table/create",
				"/api/support/table/drop",
			}

			for _, path := range pathForSkipAuthorization {
				if request.Request().URL.Path == path {
					return next(request)
				}
			}

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
