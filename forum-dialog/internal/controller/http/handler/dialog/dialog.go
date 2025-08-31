package dialog

import (
	"forum-dialog/internal/model"
	"forum-dialog/pkg/api/forum-authorization"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strconv"
)

func CreateDialog(dialogService model.IDialogService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body CreateDialogBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		dialogID, err := dialogService.CreateDialog(ctx,
			authorizationData.AccountID,
			body.Account2ID,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, CreateDialogResponse{dialogID})
	}
}

func MarkDialogAsStarred(dialogService model.IDialogService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		dialogIDStr := request.Param("dialogID")
		dialogID, err := strconv.Atoi(dialogIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		err = dialogService.MarkDialogAsStarred(ctx, authorizationData.AccountID, dialogID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, "ok")
	}
}

func UploadFile(dialogService model.IDialogService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		file, err := request.FormFile("file")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer src.Close()

		fileSrc, err := io.ReadAll(src)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		fileURL, err := dialogService.UploadFile(ctx, fileSrc, file.Filename)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, UploadFileResponse{
			FileURL: fileURL,
		})
	}
}

func MarkMessagesAsRead(dialogService model.IDialogService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		dialogIDStr := request.Param("dialogID")
		dialogID, err := strconv.Atoi(dialogIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = dialogService.MarkMessagesAsRead(ctx, dialogID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "ok")
	}
}

func DialogsByAccountID(dialogService model.IDialogService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		dialogs, err := dialogService.DialogsByAccountID(ctx, authorizationData.AccountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, DialogsByAccountIDResponse{
			Dialogs: dialogs,
		})
	}
}

func MessagesByDialogID(dialogService model.IDialogService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		dialogIDStr := request.Param("dialogID")
		dialogID, err := strconv.Atoi(dialogIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		messages, files, err := dialogService.MessagesByDialogID(ctx, dialogID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, MessagesByDialogIDResponse{
			Messages: messages,
			Files:    files,
		})
	}
}

func DeleteDialog(dialogService model.IDialogService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		dialogIDStr := request.Param("dialogID")
		dialogID, err := strconv.Atoi(dialogIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = dialogService.DeleteDialog(ctx, dialogID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, "ok")
	}
}

func DownloadFile(dialogService model.IDialogService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		fileURL := request.Param("fileURL")

		file, err := dialogService.DownloadFile(ctx, fileURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.Blob(http.StatusOK, "application/octet-stream", file)
	}
}
