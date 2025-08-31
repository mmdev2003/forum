package dialog

import (
	"forum-support/internal/model"
	forumauthorization "forum-support/pkg/api/forum-authorization"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CreateDialog(dialogService model.IDialogService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forumauthorization.AuthorizationData)

		if authorizationData.Role != model.RoleSupport {
			return echo.NewHTTPError(http.StatusForbidden)
		}

		var requestData CreateDialogRequest
		if err := request.Bind(&requestData); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		supportRequestID, err := dialogService.CreateDialog(
			ctx,
			requestData.SupportRequestID,
			requestData.UserAccountID,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, map[string]int{
			"dialogID": supportRequestID,
		})
	}
}

func MarkMessagesAsRead(dialogService model.IDialogService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forumauthorization.AuthorizationData)

		dialogIDStr := request.Param("dialogID")
		dialogID, err := strconv.Atoi(dialogIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = dialogService.MarkMessagesAsRead(ctx, dialogID, authorizationData.AccountID, authorizationData.Role)
		if err != nil {
			if err == model.ErrActionNotAllowed {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.NoContent(http.StatusOK)
	}
}

func AllDialogs(dialogService model.IDialogService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forumauthorization.AuthorizationData)

		dialogs, err := dialogService.AllDialogs(ctx, authorizationData.AccountID, authorizationData.Role)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, dialogs)
	}
}

func MessagesByDialogID(dialogService model.IDialogService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forumauthorization.AuthorizationData)

		dialogIDStr := request.Param("dialogID")
		dialogID, err := strconv.Atoi(dialogIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		messages, err := dialogService.MessagesByDialogID(ctx, dialogID, authorizationData.AccountID, authorizationData.Role)
		if err != nil {
			if err == model.ErrActionNotAllowed {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, messages)
	}
}
