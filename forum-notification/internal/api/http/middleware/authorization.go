package middleware

import (
	"forum-notification/internal/model"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthorizationMiddleware(authClient model.IAuthorizationClient) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(request echo.Context) error {
			ctx := request.Request().Context()

			pathForSkipAuthorization := []string{
				model.PREFIX + "/table/create",
				model.PREFIX + "/table/drop",
			}

			requestPath := request.Request().URL.Path
			if strings.HasPrefix(requestPath, model.PREFIX+"/create/") {
				return next(request)
			}

			for _, path := range pathForSkipAuthorization {
				if requestPath == path {
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
