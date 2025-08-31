package frame

import (
	"context"
	"forum-frame/internal/model"
	"forum-frame/pkg/api/forum-payment"
	"time"
)

func New(
	frameRepo model.IFrameRepo,
	paymentClient model.IPaymentClient,
) *ServiceStatus {
	return &ServiceStatus{
		frameRepo:     frameRepo,
		paymentClient: paymentClient,
	}
}

type ServiceStatus struct {
	frameRepo     model.IFrameRepo
	paymentClient model.IPaymentClient
}

func (s *ServiceStatus) CreatePaymentForFrame(ctx context.Context,
	frameID,
	accountID,
	duration int,
	currency string,
) (*forum_payment.PaymentData, error) {
	var amountUSD float32
	var expirationAt time.Time

	frameData, err := s.frameRepo.FrameDataByID(ctx, frameID)
	if err != nil {
		return nil, err
	}

	// expiration перенести бы на подтверждение платежа

	if duration == -1 {
		amountUSD = frameData[0].ForeverPrice
		expirationAt = time.Now().AddDate(100, 0, 0)
	} else {
		amountUSD = frameData[0].MonthlyPrice * float32(duration)
		expirationAt = time.Now().AddDate(0, duration, 0)
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

	_, err = s.frameRepo.CreateFrame(ctx, frameID, accountID, paymentData.PaymentID, expirationAt)
	if err != nil {
		return nil, err
	}

	return paymentData, nil
}

func (s *ServiceStatus) AddNewFrame(ctx context.Context,
	frameFile []byte,
	monthlyPrice,
	foreverPrice float32,
	name string,
) error {
	fileID, err := s.frameRepo.UploadFrame(frameFile, name)
	if err != nil {
		return err
	}

	err = s.frameRepo.CreateFrameData(ctx, monthlyPrice, foreverPrice, name, fileID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceStatus) ChangeCurrentFrame(ctx context.Context,
	dbFrameID,
	accountID int,
) error {
	err := s.frameRepo.ChangeCurrentFrame(ctx, dbFrameID, accountID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceStatus) ConfirmPaymentForFrame(ctx context.Context,
	paymentID int,
) error {
	err := s.frameRepo.ConfirmPaymentForFrame(ctx, paymentID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceStatus) FramesByAccountID(ctx context.Context,
	accountID int,
) ([]*model.Frame, *model.CurrentFrame, error) {
	accountFrames, err := s.frameRepo.FramesByAccountID(ctx, accountID)
	if err != nil {
		return nil, nil, err
	}
	currentFrame, err := s.frameRepo.CurrentFrameByAccountID(ctx, accountID)
	if err != nil {
		return nil, nil, err
	}

	return accountFrames, currentFrame[0], err
}

func (s *ServiceStatus) AllFrame(ctx context.Context) ([]*model.FrameData, error) {
	allFrames, err := s.frameRepo.AllFrame(ctx)
	if err != nil {
		return nil, err
	}
	return allFrames, nil
}

func (s *ServiceStatus) DownloadFrame(
	ctx context.Context,
	frameID int,
) ([]byte, error) {
	frameFile, err := s.frameRepo.DownloadFrame(ctx, frameID)
	if err != nil {
		return nil, err
	}
	return frameFile, nil
}
