package status

import (
	"context"
	"forum-status/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"time"
)

func New(db model.IDatabase) *RepoStatus {
	return &RepoStatus{db}
}

type RepoStatus struct {
	db model.IDatabase
}

func (statusRepo *RepoStatus) CreateStatus(ctx context.Context,
	statusID,
	accountID,
	paymentID int,
	expirationAt time.Time,
) (int, error) {
	args := pgx.NamedArgs{
		"status_id":      statusID,
		"account_id":     accountID,
		"payment_id":     paymentID,
		"payment_status": model.Pending,
		"expiration_at":  expirationAt,
	}
	dbStatusID, err := statusRepo.db.Insert(ctx, CreateStatus, args)
	if err != nil {
		return 0, err
	}

	return dbStatusID, err
}

func (statusRepo *RepoStatus) ConfirmPaymentForStatus(ctx context.Context,
	paymentID int,
) error {
	args := pgx.NamedArgs{
		"payment_id":     paymentID,
		"payment_status": model.Confirmed,
	}
	err := statusRepo.db.Update(ctx, ConfirmPaymentForStatus, args)
	if err != nil {
		return err
	}
	return err
}

func (statusRepo *RepoStatus) DeleteStatus(ctx context.Context,
	statusID int,
	accountID int,
) error {
	args := pgx.NamedArgs{
		"status_id":  statusID,
		"account_id": accountID,
	}
	err := statusRepo.db.Delete(ctx, DeleteStatus, args)
	if err != nil {
		return err
	}

	return nil
}

func (statusRepo *RepoStatus) StatusesByAccountID(ctx context.Context,
	accountID int,
) ([]*model.Status, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	rows, err := statusRepo.db.Select(ctx, StatusesByAccountID, args)
	if err != nil {
		return nil, err
	}

	var statuses []*model.Status
	err = pgxscan.ScanAll(&statuses, rows)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}
