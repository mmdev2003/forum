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

func CreateTopic(topicService model.ITopicService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()
		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		var body CreateTopicBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		topicOwnerAccountID := authorizationData.AccountID
		topicID, err := topicService.CreateTopic(ctx,
			body.SubthreadID,
			body.ThreadID,
			topicOwnerAccountID,
			body.ThreadName,
			body.SubthreadName,
			body.TopicName,
			body.TopicOwnerLogin,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, CreateTopicResponse{
			TopicID: topicID,
		})
	}
}

func UpdateTopicAvatar(topicService model.ITopicService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body UpdateTopicAvatarBody
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
		file := form.File["file"][0]

		f, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open file: "+err.Error())
		}
		defer f.Close()

		fileBytes, err := io.ReadAll(f)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read file: "+err.Error())
		}
		extension := filepath.Ext(file.Filename)
		fileName := strings.TrimSuffix(file.Filename, extension)

		err = topicService.UpdateTopicAvatar(ctx,
			body.TopicID,
			fileName,
			extension,
			fileBytes,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func DownloadTopicAvatar(topicService model.ITopicService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		fileURL := request.Param("fileURL")

		file, err := topicService.DownloadTopicAvatar(ctx, fileURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.Blob(http.StatusOK, "application/octet-stream", file)
	}
}

func AddViewToTopic(topicService model.ITopicService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body AddViewToTopicBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		err := topicService.AddViewToTopic(ctx, body.TopicID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func CloseTopic(topicService model.ITopicService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body CloseTopicBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		if authorizationData.Role != model.RoleAdmin {
			return echo.NewHTTPError(http.StatusForbidden, "only admin can close topic")
		}

		err := topicService.CloseTopic(ctx,
			body.TopicOwnerAccountID,
			body.AdminAccountID,
			body.TopicID,
			body.TopicName,
			body.AdminLogin,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func RejectTopic(topicService model.ITopicService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)
		if authorizationData.Role != model.RoleAdmin {
			return echo.NewHTTPError(http.StatusForbidden, "only admin can reject topic")
		}

		topicIDStr := request.Param("topicID")
		topicID, err := strconv.Atoi(topicIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		err = topicService.RejectTopic(ctx,
			topicID,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func ApproveTopic(topicService model.ITopicService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)
		if authorizationData.Role != model.RoleAdmin {
			return echo.NewHTTPError(http.StatusForbidden, "only admin can approve topic")
		}

		topicIDStr := request.Param("topicID")
		topicID, err := strconv.Atoi(topicIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		err = topicService.ApproveTopic(ctx,
			topicID,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func ChangeTopicPriority(topicService model.ITopicService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body ChangeTopicPriorityBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		err := topicService.ChangeTopicPriority(ctx,
			body.SubthreadID,
			body.TopicID,
			body.TopicPriority,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func TopicsBySubthreadID(topicService model.ITopicService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		subthreadIDStr := request.Param("subthreadID")
		subthreadID, err := strconv.Atoi(subthreadIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		topics, err := topicService.TopicsBySubthreadID(ctx, subthreadID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, TopicsBySubthreadIDResponse{
			Topics: topics,
		})
	}
}

func TopicsByAccountID(topicService model.ITopicService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		accountIDStr := request.Param("accountID")
		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		topics, err := topicService.TopicsByAccountID(ctx, accountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, TopicsByAccountIDResponse{
			Topics: topics,
		})
	}
}

func TopicsOnModeration(topicService model.ITopicService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)
		if authorizationData.Role != model.RoleAdmin {
			return echo.NewHTTPError(http.StatusForbidden, "only admin can get topics on moderation")
		}

		topics, err := topicService.TopicsOnModeration(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, TopicsOnModerationIDResponse{
			Topics: topics,
		})
	}
}
