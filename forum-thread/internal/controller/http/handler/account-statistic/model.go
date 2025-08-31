package account_statistic

import "forum-thread/internal/model"

type CreateAccountStatisticBody struct {
	AccountID int `json:"accountID" validate:"required"`
}

type StatisticByAccountIDResponse struct {
	AccountStatistic *model.AccountStatistic `json:"accountStatistic"`
}
