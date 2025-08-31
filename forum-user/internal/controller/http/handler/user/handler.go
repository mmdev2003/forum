package handler

import (
	"forum-user/internal/model"
	"forum-user/pkg/api/forum-authorization"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strconv"
)

func CreateUser(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body CreateUserBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		_, err := userService.CreateUser(ctx,
			body.AccountID,
			body.Login,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func BanUser(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body BanUserBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		err := userService.BanUser(ctx,
			authorizationData.AccountID,
			body.ToAccountID,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func NewWarningFromAdmin(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body NewWarningFromAdminBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)
		adminAccountID := authorizationData.AccountID

		if authorizationData.Role != model.RoleAdmin {
			return echo.NewHTTPError(http.StatusUnauthorized, "not authorized")
		}

		err := userService.NewWarningFromAdmin(ctx,
			adminAccountID,
			body.ToAccountID,
			body.WarningType,
			body.AdminLogin,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func UploadAvatar(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		file, err := request.FormFile("avatar")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer src.Close()

		avatar, err := io.ReadAll(src)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := userService.UploadAvatar(ctx, authorizationData.AccountID, avatar); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func UploadHeader(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		file, err := request.FormFile("header")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer src.Close()

		header, err := io.ReadAll(src)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if err := userService.UploadHeader(ctx, authorizationData.AccountID, header); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func DownloadAvatar(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		fileID := request.Param("fileID")

		avatar, err := userService.DownloadAvatar(ctx, fileID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.Blob(http.StatusOK, "application/octet-stream", avatar)
	}
}

func DownloadHeader(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		fileID := request.Param("fileID")

		header, err := userService.DownloadHeader(ctx, fileID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.Blob(http.StatusOK, "application/octet-stream", header)
	}
}

func GetUserByAccountID(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		accountIDStr := request.Param("accountID")
		accountID, _ := strconv.Atoi(accountIDStr)

		user, err := userService.UserByAccountID(ctx, accountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, GetUserByAccountIDResponse{
			AccountID: user.AccountID,
			Login:     user.Login,
			HeaderUrl: user.HeaderUrl,
			AvatarUrl: user.AvatarUrl,
			CreatedAt: user.CreatedAt,
		})
	}
}

func GetUserByLogin(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		login := request.Param("login")

		user, err := userService.UserByLogin(ctx, login)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, GetUserByLoginResponse{
			User: user,
		})
	}
}

func UsersByLoginSearch(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		login := request.Param("login")

		users, err := userService.UsersByLoginSearch(ctx, login)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, UsersByLoginResponse{
			Users: users,
		})
	}
}

func BanByAccountID(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		userBans, err := userService.BanByAccountID(ctx, authorizationData.AccountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, BanByAccountIDResponse{
			UserBans: userBans,
		})
	}
}

func AllWarningFromAdmin(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		toAccountIDStr := request.Param("toAccountID")
		toAccountID, _ := strconv.Atoi(toAccountIDStr)

		userWarnings, err := userService.AllWarningFromAdmin(ctx, toAccountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, AllWarningFromAdminResponse{
			UserWarnings: userWarnings,
		})
	}
}

func UnbanUser(userService model.IUserService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		toAccountIDStr := request.Param("toAccountID")
		toAccountID, err := strconv.Atoi(toAccountIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		err = userService.UnbanUser(ctx, authorizationData.AccountID, toAccountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, "ok")
	}
}
