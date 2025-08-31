package support_request

import (
	"context"
	"forum-support/internal/model"
)

func New(
	dialogRepo model.ISupportRequestRepo,
	notificationClient model.INotificationClient,
) *SupportRequestService {
	return &SupportRequestService{
		supportRequestRepo: dialogRepo,
		notificationClient: notificationClient,
	}
}

type SupportRequestService struct {
	supportRequestRepo model.ISupportRequestRepo
	notificationClient model.INotificationClient
}

func (s *SupportRequestService) GetRequestById(ctx context.Context, supportRequestID int) (*model.SupportRequest, error) {
	return s.supportRequestRepo.GetRequestById(ctx, supportRequestID)
}

func (s *SupportRequestService) GetRequests(ctx context.Context) ([]model.SupportRequest, error) {
	return s.supportRequestRepo.GetRequests(ctx)
}

func (s *SupportRequestService) GetRequestsWithStatus(ctx context.Context, status model.RequestStatus) ([]model.SupportRequest, error) {
	return s.supportRequestRepo.GetRequestsWithStatus(ctx, status)
}

func (s SupportRequestService) CreateRequest(ctx context.Context, accountID int, title, description string) (int, error) {
	return s.supportRequestRepo.CreateRequest(ctx, accountID, title, description)
}

func (s SupportRequestService) OpenRequest(ctx context.Context, supportRequestID int, role string) error {
	request, err := s.supportRequestRepo.GetRequestById(ctx, supportRequestID)
	if err != nil {
		return err
	}

	if role != model.RoleSupport {
		return model.ErrActionNotAllowed
	}

	err = s.supportRequestRepo.OpenRequest(ctx, supportRequestID)
	if err != nil {
		return err
	}

	return s.notificationClient.SendStatusReceivedNotificationRequest(request.AccountID, model.OpenRequestStatus)
}

func (s SupportRequestService) CloseRequest(ctx context.Context, supportRequestID, actorAccountID int, role string) error {
	request, err := s.supportRequestRepo.GetRequestById(ctx, supportRequestID)
	if err != nil {
		return err
	}

	if role != model.RoleSupport && request.AccountID != actorAccountID {
		return model.ErrActionNotAllowed
	}

	err = s.supportRequestRepo.CloseRequest(ctx, supportRequestID)
	if err != nil {
		return err
	}

	return s.notificationClient.SendSupportRequestClosedNotification(request.AccountID, supportRequestID)
}
