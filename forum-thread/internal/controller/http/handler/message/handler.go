package message

import (
	"encoding/json"
	"forum-thread/internal/model"
	"forum-thread/pkg/api/forum-authorization"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func SendMessageToTopic(messageService model.IMessageService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		var body SendMessageToTopicBody
		if err := json.Unmarshal([]byte(request.FormValue("body")), &body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		form, err := request.MultipartForm()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse multipart form: "+err.Error())
		}
		files := form.File["files"]

		var uploadedFiles [][]byte
		var filesNames []string
		var filesExtensions []string
		if files != nil {
			for _, file := range files {
				f, err := file.Open()
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file: "+err.Error())
				}
				defer f.Close()

				bytes, err := io.ReadAll(f)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read file: "+err.Error())
				}

				uploadedFiles = append(uploadedFiles, bytes)
				extension := filepath.Ext(file.Filename)

				fileName := strings.TrimSuffix(file.Filename, extension)
				filesNames = append(filesNames, fileName)
				filesExtensions = append(filesExtensions, extension)
			}
		}

		senderAccountID := authorizationData.AccountID

		filesURLs, err := messageService.SendMessageToTopic(ctx,
			body.SubthreadID,
			body.TopicID,
			body.ReplyToMessageID,
			body.ReplyMessageOwnerAccountID,
			body.TopicOwnerAccountID,
			senderAccountID,
			body.SenderLogin,
			body.ThreadName,
			body.SubthreadName,
			body.TopicName,
			body.SenderMessageText,
			uploadedFiles,
			filesNames,
			filesExtensions,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, SendMessageToTopicResponse{
			FilesURLs: filesURLs,
		})
	}
}

func LikeMessage(messageService model.IMessageService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body LikeMessageBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		err := messageService.LikeMessage(ctx,
			body.TopicID,
			body.MessageOwnerAccountID,
			body.LikerAccountID,
			body.LikeMessageID,
			body.LikeTypeID,
			body.LikerLogin,
			body.TopicName,
			body.LikeMessageText,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func UnlikeMessage(messageService model.IMessageService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body UnlikeMessageBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		likerAccountID := authorizationData.AccountID
		err := messageService.UnlikeMessage(ctx,
			body.LikeMessageID,
			likerAccountID,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func ReportMessage(messageService model.IMessageService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body ReportMessageBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		err := messageService.ReportMessage(ctx,
			body.MessageID,
			body.AccountID,
			body.ReportText,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func MessagesByTopicID(messageService model.IMessageService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		topicIDStr := request.Param("topicID")
		topicID, err := strconv.Atoi(topicIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		messages, likes, files, err := messageService.MessagesByTopicID(ctx,
			authorizationData.AccountID,
			topicID,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, MessagesByTopicIDResponse{
			Messages: messages,
			Likes:    likes,
			Files:    files,
		})
	}
}

func MessagesByAccountID(messageService model.IMessageService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		accountIDStr := request.Param("accountID")
		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		messages, files, err := messageService.MessagesByAccountID(ctx,
			accountID,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, MessagesByAccountIDResponse{
			Messages: messages,
			Files:    files,
		})
	}
}

func EditMessage(messageService model.IMessageService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body EditMessageBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		err := messageService.EditMessage(ctx, body.MessageID, body.MessageText)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func MessagesByText(messageService model.IMessageService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		text := request.Param("text")

		messages, err := messageService.MessagesByText(ctx, text)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, MessagesByTextResponse{
			Messages: messages,
		})
	}
}

func DownloadFile(messageService model.IMessageService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		fileURL := request.Param("fileURL")

		file, err := messageService.DownloadFile(ctx, fileURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.Blob(http.StatusOK, "application/octet-stream", file)
	}
}
