package filter

import (
	"context"
	"forum-notification/internal/model"
	"log/slog"
)

func New(
	notificationFilterRepo model.INotificationFilterRepo,
) *ServiceNotificationFilter {
	return &ServiceNotificationFilter{
		notificationFilterRepo: notificationFilterRepo,
	}
}

type ServiceNotificationFilter struct {
	notificationFilterRepo model.INotificationFilterRepo
}

func (s *ServiceNotificationFilter) IsNotificationEnabled(ctx context.Context, accountID int, notificationType model.NotificationType) (bool, error) {
	return s.notificationFilterRepo.IsNotificationEnabled(ctx, accountID, notificationType)
}

func (s *ServiceNotificationFilter) GetFilters(ctx context.Context, accountID int) ([]model.NotificationType, error) {
	filters, dberr := s.notificationFilterRepo.GetFilters(ctx, accountID)
	slog.Error("Error while getting filters: ", "message", dberr)

	if filters == nil {
		filters = []model.NotificationType{
			model.MessageFromTopicType,
			model.MessageReplyFromTopicType,
			model.LikeMessageFromTopicType,
			model.TopicClosedType,
			model.ResponseToSupportRequestType,
			model.StatusReceivedType,
			model.FrameReceivedType,
			model.MessageFromDialogType,
			model.MentionFromTopicType,
			model.WarningFromAdminType,
		}
		err := s.UpsertFilters(ctx, accountID, filters)
		if err != nil {
			return nil, err
		}
	}

	return filters, dberr
}

func (s *ServiceNotificationFilter) UpsertFilters(ctx context.Context, accountId int, settings []model.NotificationType) error {
	return s.notificationFilterRepo.UpsertFilters(ctx, accountId, settings)
}
