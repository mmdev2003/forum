package notification

import (
	"forum-notification/internal/model"
	forumauthorization "forum-notification/pkg/api/forum-authorization"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateNotification(notificationService model.INotification) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()
		notificationType := model.NotificationType(request.Param("type"))

		handler, exists := notificationHandlers[notificationType]

		if !exists {
			return echo.NewHTTPError(http.StatusBadRequest, "unsupported notification type")
		}

		body, err := handler.Bind(request)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		notificationID, err := handler.Process(ctx, notificationService, body)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if notificationID == 0 {
			return request.NoContent(http.StatusNoContent)
		}

		return request.JSON(http.StatusCreated, CreateNotificationResponse{notificationID})
	}
}

func GetNotificationsByAccountID(notificationService model.INotification) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forumauthorization.AuthorizationData)

		notifications, err := notificationService.GetNotificationsByAccountID(ctx, authorizationData.AccountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, notifications)
	}
}
