package notification

import (
	"context"
	"encoding/json"
	"forum-notification/internal/model"
	"time"
)

func New(
	notificationRepo model.INotificationRepo,
	notificationFilter model.INotificationFilter,
	wsConnManager model.IWsConnManager,
) *ServiceNotification {
	return &ServiceNotification{
		notificationRepo:   notificationRepo,
		notificationFilter: notificationFilter,
		wsConnManager:      wsConnManager,
	}
}

type ServiceNotification struct {
	notificationRepo   model.INotificationRepo
	notificationFilter model.INotificationFilter
	wsConnManager      model.IWsConnManager
}

func (s *ServiceNotification) CreateMessageFromTopicNotification(ctx context.Context,
	accountID,
	messageID,
	replierAccountID,
	topicID int,
	messageText,
	topicName,
	replierLogin string,
) (int, error) {
	isNotificationEnabled, _ := s.notificationFilter.IsNotificationEnabled(ctx, accountID, model.MessageFromTopicType)
	if isNotificationEnabled {
		return 0, nil
	}

	notificationID, err := s.notificationRepo.CreateMessageFromTopicNotification(ctx,
		accountID,
		messageID,
		replierAccountID,
		topicID,
		messageText,
		topicName,
		replierLogin,
	)
	if err != nil {
		return 0, err
	}

	payload, _ := json.Marshal(MessageFromTopicNotificationPayload{
		MessageID:        messageID,
		ReplierAccountID: replierAccountID,
		TopicID:          topicID,
		MessageText:      messageText,
		TopicName:        topicName,
		ReplierLogin:     replierLogin,
	})
	wsBody, _ := json.Marshal(NotificationWsBody{
		ID:        notificationID,
		CreatedAt: time.Now(),
		Type:      model.MessageFromTopicType,
		Payload:   payload,
	})
	err = s.wsConnManager.SendMsg(accountID, wsBody)
	if err != nil {
		return 0, err
	}
	return notificationID, nil
}

func (s *ServiceNotification) CreateMessageReplyFromTopicNotification(ctx context.Context,
	accountID,
	messageID,
	replierAccountID,
	topicID int,
	messageText,
	topicName,
	replierLogin string,
) (int, error) {
	isNotificationEnabled, _ := s.notificationFilter.IsNotificationEnabled(ctx, accountID, model.MessageReplyFromTopicType)
	if isNotificationEnabled {
		return 0, nil
	}

	notificationID, err := s.notificationRepo.CreateMessageReplyFromTopicNotification(ctx,
		accountID,
		messageID,
		replierAccountID,
		topicID,
		messageText,
		topicName,
		replierLogin,
	)
	if err != nil {
		return 0, err
	}

	payload, _ := json.Marshal(NewMessageReplyFromTopicNotificationPayload{
		MessageID:        messageID,
		ReplierAccountID: replierAccountID,
		TopicID:          topicID,
		MessageText:      messageText,
		TopicName:        topicName,
		ReplierLogin:     replierLogin,
	})
	wsBody, _ := json.Marshal(NotificationWsBody{
		ID:        notificationID,
		CreatedAt: time.Now(),
		Type:      model.MessageReplyFromTopicType,
		Payload:   payload,
	})
	err = s.wsConnManager.SendMsg(accountID, wsBody)
	if err != nil {
		return 0, err
	}
	return notificationID, nil
}

func (s *ServiceNotification) CreateLikeMessageFromTopicNotification(ctx context.Context,
	accountID,
	messageID,
	likerAccountID,
	topicID int,
	messageText,
	topicName,
	likerLogin string,
) (int, error) {
	isNotificationEnabled, _ := s.notificationFilter.IsNotificationEnabled(ctx, accountID, model.LikeMessageFromTopicType)
	if isNotificationEnabled {
		return 0, nil
	}

	notificationID, err := s.notificationRepo.CreateLikeMessageFromTopicNotification(ctx,
		accountID,
		messageID,
		likerAccountID,
		topicID,
		messageText,
		topicName,
		likerLogin,
	)
	if err != nil {
		return 0, err
	}

	payload, _ := json.Marshal(NewLikeMessageFromTopicNotificationPayload{
		MessageID:      messageID,
		LikerAccountID: likerAccountID,
		TopicID:        topicID,
		MessageText:    messageText,
		TopicName:      topicName,
		LikerLogin:     likerLogin,
	})
	wsBody, _ := json.Marshal(NotificationWsBody{
		ID:        notificationID,
		CreatedAt: time.Now(),
		Type:      model.LikeMessageFromTopicType,
		Payload:   payload,
	})
	err = s.wsConnManager.SendMsg(accountID, wsBody)
	if err != nil {
		return 0, err
	}
	return notificationID, nil
}

func (s *ServiceNotification) CreateTopicClosedNotification(ctx context.Context,
	accountID,
	adminAccountID,
	topicID int,
	topicName,
	adminLogin string,
) (int, error) {
	isNotificationEnabled, _ := s.notificationFilter.IsNotificationEnabled(ctx, accountID, model.TopicClosedType)
	if isNotificationEnabled {
		return 0, nil
	}

	notificationID, err := s.notificationRepo.CreateTopicClosedNotification(ctx,
		accountID,
		adminAccountID,
		topicID,
		topicName,
		adminLogin,
	)
	if err != nil {
		return 0, err
	}

	payload, _ := json.Marshal(TopicClosedNotificationPayload{
		AdminAccountID: adminAccountID,
		TopicID:        topicID,
		TopicName:      topicName,
		AdminLogin:     adminLogin,
	})
	wsBody, _ := json.Marshal(NotificationWsBody{
		ID:        notificationID,
		CreatedAt: time.Now(),
		Type:      model.TopicClosedType,
		Payload:   payload,
	})
	err = s.wsConnManager.SendMsg(accountID, wsBody)
	if err != nil {
		return 0, err
	}
	return notificationID, nil
}

func (s *ServiceNotification) CreateResponseToSupportRequestNotification(ctx context.Context,
	accountID,
	supportRequestID int,
) (int, error) {
	isNotificationEnabled, _ := s.notificationFilter.IsNotificationEnabled(ctx, accountID, model.ResponseToSupportRequestType)
	if isNotificationEnabled {
		return 0, nil
	}

	notificationID, err := s.notificationRepo.CreateResponseToSupportRequestNotification(ctx,
		accountID,
		supportRequestID,
	)
	if err != nil {
		return 0, err
	}

	payload, _ := json.Marshal(ResponseToSupportRequestNotificationPayload{
		SupportRequestID: supportRequestID,
	})
	wsBody, _ := json.Marshal(NotificationWsBody{
		ID:        notificationID,
		CreatedAt: time.Now(),
		Type:      model.ResponseToSupportRequestType,
		Payload:   payload,
	})
	err = s.wsConnManager.SendMsg(accountID, wsBody)
	if err != nil {
		return 0, err
	}
	return notificationID, nil
}

func (s *ServiceNotification) CreateStatusReceivedNotification(ctx context.Context,
	accountID int,
	statusName string,
) (int, error) {
	isNotificationEnabled, _ := s.notificationFilter.IsNotificationEnabled(ctx, accountID, model.StatusReceivedType)
	if isNotificationEnabled {
		return 0, nil
	}

	notificationID, err := s.notificationRepo.CreateStatusReceivedNotification(ctx,
		accountID,
		statusName,
	)
	if err != nil {
		return 0, err
	}

	payload, _ := json.Marshal(StatusReceivedNotificationPayload{
		StatusName: statusName,
	})
	wsBody, _ := json.Marshal(NotificationWsBody{
		ID:        notificationID,
		CreatedAt: time.Now(),
		Type:      model.StatusReceivedType,
		Payload:   payload,
	})
	err = s.wsConnManager.SendMsg(accountID, wsBody)
	if err != nil {
		return 0, err
	}
	return notificationID, nil
}

func (s *ServiceNotification) CreateFrameReceivedNotification(ctx context.Context,
	accountID int,
	frameName string,
) (int, error) {
	isNotificationEnabled, _ := s.notificationFilter.IsNotificationEnabled(ctx, accountID, model.FrameReceivedType)
	if isNotificationEnabled {
		return 0, nil
	}

	notificationID, err := s.notificationRepo.CreateFrameReceivedNotification(ctx,
		accountID,
		frameName,
	)
	if err != nil {
		return 0, err
	}

	payload, _ := json.Marshal(FrameReceivedNotificationPayload{
		FrameName: frameName,
	})
	wsBody, _ := json.Marshal(NotificationWsBody{
		ID:        notificationID,
		CreatedAt: time.Now(),
		Type:      model.FrameReceivedType,
		Payload:   payload,
	})
	err = s.wsConnManager.SendMsg(accountID, wsBody)
	if err != nil {
		return 0, err
	}
	return notificationID, nil
}

func (s *ServiceNotification) CreateMessageFromDialogNotification(ctx context.Context,
	accountID,
	messageID,
	dialogID,
	senderAccountID int,
	messageText,
	senderLogin string,
) (int, error) {
	isNotificationEnabled, _ := s.notificationFilter.IsNotificationEnabled(ctx, accountID, model.MessageFromDialogType)
	if isNotificationEnabled {
		return 0, nil
	}

	notificationID, err := s.notificationRepo.CreateMessageFromDialogNotification(ctx,
		accountID,
		messageID,
		dialogID,
		senderAccountID,
		messageText,
		senderLogin,
	)
	if err != nil {
		return 0, err
	}

	payload, _ := json.Marshal(NewMessageFromDialogNotificationPayload{
		MessageID:       messageID,
		DialogID:        dialogID,
		SenderAccountID: senderAccountID,
		MessageText:     messageText,
		SenderLogin:     senderLogin,
	})
	wsBody, _ := json.Marshal(NotificationWsBody{
		ID:        notificationID,
		CreatedAt: time.Now(),
		Type:      model.MessageFromDialogType,
		Payload:   payload,
	})
	err = s.wsConnManager.SendMsg(accountID, wsBody)
	if err != nil {
		return 0, err
	}
	return notificationID, nil

}

func (s *ServiceNotification) CreateMentionFromTopicNotification(ctx context.Context,
	accountID,
	messageID,
	mentionAccountID int,
	messageText,
	topicName,
	mentionLogin string,
) (int, error) {
	isNotificationEnabled, _ := s.notificationFilter.IsNotificationEnabled(ctx, accountID, model.MentionFromTopicType)
	if isNotificationEnabled {
		return 0, nil
	}

	notificationID, err := s.notificationRepo.CreateMentionFromTopicNotification(ctx,
		accountID,
		messageID,
		mentionAccountID,
		messageText,
		topicName,
		mentionLogin,
	)
	if err != nil {
		return 0, err
	}

	payload, _ := json.Marshal(NewMentionFromTopicNotificationPayload{
		MessageID:        messageID,
		MentionAccountID: mentionAccountID,
		MessageText:      messageText,
		TopicName:        topicName,
		MentionLogin:     mentionLogin,
	})
	wsBody, _ := json.Marshal(NotificationWsBody{
		ID:        notificationID,
		CreatedAt: time.Now(),
		Type:      model.MentionFromTopicType,
		Payload:   payload,
	})
	err = s.wsConnManager.SendMsg(accountID, wsBody)
	if err != nil {
		return 0, err
	}
	return notificationID, nil
}

func (s *ServiceNotification) CreateWarningFromAdminNotification(ctx context.Context,
	accountID,
	adminAccountID int,
	warningText,
	adminLogin string,
) (int, error) {
	isNotificationEnabled, _ := s.notificationFilter.IsNotificationEnabled(ctx, accountID, model.WarningFromAdminType)
	if isNotificationEnabled {
		return 0, nil
	}

	notificationID, err := s.notificationRepo.CreateWarningFromAdminNotification(ctx,
		accountID,
		adminAccountID,
		warningText,
		adminLogin,
	)
	if err != nil {
		return 0, err
	}

	payload, _ := json.Marshal(NewWarningFromAdminNotificationPayload{
		AdminAccountID: adminAccountID,
		WarningText:    warningText,
		AdminLogin:     adminLogin,
	})
	wsBody, _ := json.Marshal(NotificationWsBody{
		ID:        notificationID,
		CreatedAt: time.Now(),
		Type:      model.WarningFromAdminType,
		Payload:   payload,
	})
	err = s.wsConnManager.SendMsg(accountID, wsBody)
	if err != nil {
		return 0, err
	}
	return notificationID, nil
}

func (s *ServiceNotification) GetNotificationsByAccountID(ctx context.Context, id int) ([]interface{}, error) {
	return s.notificationRepo.GetNotificationsByAccountID(ctx, id)
}
