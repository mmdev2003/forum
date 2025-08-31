package status

import (
	"context"
	"forum-status/internal/model"
	"forum-status/pkg/api/forum-payment"
	"time"
)

func New(
	statusRepo model.IStatusRepo,
	paymentClient model.IPaymentClient,
) *ServiceStatus {
	return &ServiceStatus{
		statusRepo:    statusRepo,
		paymentClient: paymentClient,
	}
}

type ServiceStatus struct {
	statusRepo    model.IStatusRepo
	paymentClient model.IPaymentClient
}

func (s *ServiceStatus) CreatePaymentForStatus(ctx context.Context,
	statusID,
	accountID,
	duration int,
	currency string,
) (*forum_payment.PaymentData, error) {
	var amountUSD float32
	var expirationAt time.Time

	// expiration перенести бы на подтверждение платежа
	for _, status := range model.Statuses {
		if status.ID == statusID {
			if duration == -1 {
				amountUSD = status.ForeverPrice
				expirationAt = time.Now().AddDate(100, 0, 0)
			} else {
				amountUSD = status.MonthlyPrice * float32(duration)
				expirationAt = time.Now().AddDate(0, duration, 0)
			}
		}
	}

	paymentData, err := s.paymentClient.CreatePayment(ctx,
		accountID,
		string(model.TStatus),
		currency,
		amountUSD,
	)
	if err != nil {
		return nil, err
	}

	_, err = s.statusRepo.CreateStatus(ctx, statusID, accountID, paymentData.PaymentID, expirationAt)
	if err != nil {
		return nil, err
	}

	return paymentData, nil

}

func (s *ServiceStatus) AssignStatusToAccount(ctx context.Context,
	statusID,
	accountID int,
) error {
	expirationAt := time.Now().AddDate(100, 0, 0)
	_, err := s.statusRepo.CreateStatus(ctx, statusID, accountID, 0, expirationAt)
	if err != nil {
		return err
	}

	err = s.statusRepo.ConfirmPaymentForStatus(ctx, 0)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceStatus) RevokeStatus(ctx context.Context,
	statusID,
	accountID int,
) error {
	err := s.statusRepo.DeleteStatus(ctx, statusID, accountID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceStatus) ConfirmPaymentForStatus(ctx context.Context,
	paymentID int,
) error {
	err := s.statusRepo.ConfirmPaymentForStatus(ctx, paymentID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceStatus) StatusesByAccountID(ctx context.Context,
	accountID int,
) ([]*model.Status, error) {
	accountStatuses, err := s.statusRepo.StatusesByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return accountStatuses, nil
}
