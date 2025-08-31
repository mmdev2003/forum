package notification

import (
	"context"
	"forum-notification/internal/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateNotificationHandler interface {
	Bind(request echo.Context) (interface{}, error)
	Process(ctx context.Context, notificationService model.INotification, notification interface{}) (int, error)
}

var notificationHandlers = map[model.NotificationType]CreateNotificationHandler{
	model.MessageFromTopicType:         &MessageFromTopicHandler{},
	model.MessageReplyFromTopicType:    &MessageReplyFromTopicHandler{},
	model.LikeMessageFromTopicType:     &LikeMessageFromTopicHandler{},
	model.TopicClosedType:              &TopicClosedHandler{},
	model.ResponseToSupportRequestType: &ResponseToSupportRequestHandler{},
	model.StatusReceivedType:           &StatusReceivedHandler{},
	model.FrameReceivedType:            &FrameReceivedHandler{},
	model.MessageFromDialogType:        &MessageFromDialogHandler{},
	model.MentionFromTopicType:         &MentionFromTopicHandler{},
	model.WarningFromAdminType:         &WarningFromAdminHandler{},
}

type MessageFromTopicHandler struct{}

func (m MessageFromTopicHandler) Bind(request echo.Context) (interface{}, error) {
	var notification NewMessageFromTopicNotificationRequest
	if err := request.Bind(&notification); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return notification, nil
}

func (m MessageFromTopicHandler) Process(ctx context.Context, notificationService model.INotification, payload interface{}) (int, error) {
	notification := payload.(NewMessageFromTopicNotificationRequest)
	return notificationService.CreateMessageFromTopicNotification(ctx,
		notification.TopicOwnerAccountID,
		notification.SenderMessageID,
		notification.SenderAccountID,
		notification.TopicID,
		notification.SenderMessageText,
		notification.TopicName,
		notification.SenderLogin,
	)
}

type MessageReplyFromTopicHandler struct{}

func (m MessageReplyFromTopicHandler) Bind(request echo.Context) (interface{}, error) {
	var notification NewMessageReplyFromTopicNotificationRequest
	if err := request.Bind(&notification); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return notification, nil
}

func (m MessageReplyFromTopicHandler) Process(ctx context.Context, notificationService model.INotification, payload interface{}) (int, error) {
	notification := payload.(NewMessageReplyFromTopicNotificationRequest)
	return notificationService.CreateMessageReplyFromTopicNotification(ctx,
		notification.ReplyMessageOwnerAccountID,
		notification.SenderMessageID,
		notification.SenderAccountID,
		notification.TopicID,
		notification.SenderMessageText,
		notification.TopicName,
		notification.SenderLogin,
	)
}

type LikeMessageFromTopicHandler struct{}

func (l LikeMessageFromTopicHandler) Bind(request echo.Context) (interface{}, error) {
	var notification NewLikeMessageFromTopicNotificationRequest
	if err := request.Bind(&notification); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return notification, nil
}

func (l LikeMessageFromTopicHandler) Process(ctx context.Context, notificationService model.INotification, payload interface{}) (int, error) {
	notification := payload.(NewLikeMessageFromTopicNotificationRequest)
	return notificationService.CreateLikeMessageFromTopicNotification(ctx,
		notification.MessageOwnerAccountID,
		notification.SenderMessageID,
		notification.LikerAccountID,
		notification.TopicID,
		notification.SenderMessageText,
		notification.TopicName,
		notification.LikerLogin,
	)
}

type TopicClosedHandler struct{}

func (t TopicClosedHandler) Bind(request echo.Context) (interface{}, error) {
	var notification TopicClosedNotificationRequest
	if err := request.Bind(&notification); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return notification, nil
}

func (t TopicClosedHandler) Process(ctx context.Context, notificationService model.INotification, payload interface{}) (int, error) {
	notification := payload.(TopicClosedNotificationRequest)
	return notificationService.CreateTopicClosedNotification(ctx,
		notification.TopicOwnerAccountID,
		notification.AdminAccountID,
		notification.TopicID,
		notification.TopicName,
		notification.AdminLogin,
	)
}

type ResponseToSupportRequestHandler struct{}

func (r ResponseToSupportRequestHandler) Bind(request echo.Context) (interface{}, error) {
	var notification ResponseToSupportRequestNotificationRequest
	if err := request.Bind(&notification); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return notification, nil
}

func (r ResponseToSupportRequestHandler) Process(ctx context.Context, notificationService model.INotification, payload interface{}) (int, error) {
	notification := payload.(ResponseToSupportRequestNotificationRequest)
	return notificationService.CreateResponseToSupportRequestNotification(ctx,
		notification.RequesterAccountID,
		notification.SupportRequestID,
	)
}

type StatusReceivedHandler struct{}

func (s StatusReceivedHandler) Bind(request echo.Context) (interface{}, error) {
	var notification StatusReceivedNotificationRequest
	if err := request.Bind(&notification); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return notification, nil
}

func (s StatusReceivedHandler) Process(ctx context.Context, notificationService model.INotification, payload interface{}) (int, error) {
	notification := payload.(StatusReceivedNotificationRequest)
	return notificationService.CreateStatusReceivedNotification(ctx,
		notification.ReceiverAccountID,
		notification.StatusName,
	)
}

type FrameReceivedHandler struct{}

func (f FrameReceivedHandler) Bind(request echo.Context) (interface{}, error) {
	var notification FrameReceivedNotificationRequest
	if err := request.Bind(&notification); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return notification, nil
}

func (f FrameReceivedHandler) Process(ctx context.Context, notificationService model.INotification, payload interface{}) (int, error) {
	notification := payload.(FrameReceivedNotificationRequest)
	return notificationService.CreateFrameReceivedNotification(ctx,
		notification.ReceiverAccountID,
		notification.FrameName,
	)
}

type MessageFromDialogHandler struct{}

func (m MessageFromDialogHandler) Bind(request echo.Context) (interface{}, error) {
	var notification NewMessageFromDialogNotificationRequest
	if err := request.Bind(&notification); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return notification, nil
}

func (m MessageFromDialogHandler) Process(ctx context.Context, notificationService model.INotification, payload interface{}) (int, error) {
	notification := payload.(NewMessageFromDialogNotificationRequest)
	return notificationService.CreateMessageFromDialogNotification(ctx,
		notification.AccountID,
		notification.MessageID,
		notification.DialogID,
		notification.SenderAccountID,
		notification.MessageText,
		notification.SenderLogin,
	)
}

type MentionFromTopicHandler struct{}

func (m MentionFromTopicHandler) Bind(request echo.Context) (interface{}, error) {
	var notification NewMentionFromTopicNotificationRequest
	if err := request.Bind(&notification); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return notification, nil
}

func (m MentionFromTopicHandler) Process(ctx context.Context, notificationService model.INotification, payload interface{}) (int, error) {
	notification := payload.(NewMentionFromTopicNotificationRequest)
	return notificationService.CreateMentionFromTopicNotification(ctx,
		notification.MentionedAccountID,
		notification.MessageID,
		notification.SenderAccountID,
		notification.SenderMessageText,
		notification.TopicName,
		notification.SenderLogin,
	)
}

type WarningFromAdminHandler struct{}

func (w WarningFromAdminHandler) Bind(request echo.Context) (interface{}, error) {
	var notification NewWarningFromAdminNotificationRequest
	if err := request.Bind(&notification); err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return notification, nil
}

func (w WarningFromAdminHandler) Process(ctx context.Context, notificationService model.INotification, payload interface{}) (int, error) {
	notification := payload.(NewWarningFromAdminNotificationRequest)
	return notificationService.CreateWarningFromAdminNotification(ctx,
		notification.AccountID,
		notification.AdminAccountID,
		notification.WarningText,
		notification.AdminLogin,
	)
}
