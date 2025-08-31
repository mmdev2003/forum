package forum_status

import (
	"context"
	"forum-thread/pkg/api/forum-status"
	"time"
)

func New() *StatusClient {
	return &StatusClient{}
}

type StatusClient struct{}

func (c *StatusClient) StatusByAccountID(ctx context.Context,
	accountID int,
) ([]*forum_status.Status, error) {
	if accountID == 1 {
		return []*forum_status.Status{
			{
				ID:            1,
				StatusID:      1,
				AccountID:     1,
				PaymentID:     1,
				PaymentStatus: "PAID",

				ExpirationAt: time.Time{},
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			},
			{
				ID:            2,
				StatusID:      2,
				AccountID:     1,
				PaymentID:     2,
				PaymentStatus: "PAID",

				ExpirationAt: time.Time{},
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			},
		}, nil
	}
	if accountID == 2 {
		return []*forum_status.Status{
			{
				ID:            3,
				StatusID:      2,
				AccountID:     2,
				PaymentID:     2,
				PaymentStatus: "PAID",

				ExpirationAt: time.Time{},
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			},
		}, nil
	}
	if accountID == 3 {
		return []*forum_status.Status{
			{
				ID:            4,
				StatusID:      3,
				AccountID:     3,
				PaymentID:     4,
				PaymentStatus: "PAID",

				ExpirationAt: time.Time{},
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			},
		}, nil
	}
	if accountID == 4 {
		return []*forum_status.Status{
			{
				ID:            5,
				StatusID:      4,
				AccountID:     4,
				PaymentID:     5,
				PaymentStatus: "PAID",

				ExpirationAt: time.Time{},
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			},
		}, nil
	}
	return []*forum_status.Status{}, nil
}

func (c *StatusClient) AssignStatusToAccount(ctx context.Context,
	statusID, accountID int,
) error {
	return nil
}

func (c *StatusClient) RevokeStatusFromAccount(ctx context.Context,
	statusID, accountID int,
) error {
	return nil
}
