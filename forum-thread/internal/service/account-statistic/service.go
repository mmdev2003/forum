package account_statistic

import (
	"context"
	"forum-thread/internal/model"
)

func New(
	accountStatisticRepo model.IAccountStatisticRepo,
) *ServiceAccountStatistic {
	return &ServiceAccountStatistic{
		accountStatisticRepo,
	}
}

type ServiceAccountStatistic struct {
	accountStatisticRepo model.IAccountStatisticRepo
}

func (s *ServiceAccountStatistic) CreateAccountStatistic(ctx context.Context,
	accountID int,
) (int, error) {
	accountStatisticID, err := s.accountStatisticRepo.CreateAccountStatistic(ctx, accountID)
	if err != nil {
		return 0, err
	}

	return accountStatisticID, nil
}
func (s *ServiceAccountStatistic) StatisticByAccountID(ctx context.Context,
	accountID int,
) (*model.AccountStatistic, error) {
	accountStatistic, err := s.accountStatisticRepo.StatisticByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if len(accountStatistic) == 0 {
		return nil, model.ErrAccountStatisticNotFound
	}

	return accountStatistic[0], nil
}
