package admin

import (
	"forum-admin/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateAdmin(adminService model.IAdminService, interServerSecretKey string) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body CreateAdminBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if body.InterServerSecretKey != interServerSecretKey {
			return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect secret key")
		}

		_, err := adminService.CreateAdmin(ctx,
			body.AccountID,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "ok")
	}
}

func AllAdmin(adminService model.IAdminService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		admins, err := adminService.AllAdmin(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, AllAdminResponse{
			Admins: admins,
		})
	}
}
