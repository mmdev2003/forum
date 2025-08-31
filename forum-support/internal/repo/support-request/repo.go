package support_request

import (
	"context"
	"forum-support/internal/model"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func New(db model.IDatabase) *SupportRequestRepo {
	return &SupportRequestRepo{db: db}
}

type SupportRequestRepo struct {
	db model.IDatabase
}

func (r SupportRequestRepo) GetRequestById(ctx context.Context, supportRequestID int) (*model.SupportRequest, error) {
	args := pgx.NamedArgs{
		"request_id": supportRequestID,
	}
	rows, err := r.db.Select(ctx, GetRequestById, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, pgx.ErrNoRows
	}

	var request model.SupportRequest
	err = rows.Scan(
		&request.ID,
		&request.AccountID,
		&request.Title,
		&request.Description,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r SupportRequestRepo) CreateRequest(ctx context.Context, accountID int, title, description string) (int, error) {
	args := pgx.NamedArgs{
		"account_id":  accountID,
		"title":       title,
		"description": description,
	}
	requestID, err := r.db.Insert(ctx, CreateSupportRequest, args)
	if err != nil {
		return 0, err
	}

	return requestID, nil
}

func (r SupportRequestRepo) OpenRequest(ctx context.Context, supportRequestID int) error {
	return r.setRequestStatus(ctx, supportRequestID, model.OpenRequestStatus)
}

func (r SupportRequestRepo) CloseRequest(ctx context.Context, supportRequestID int) error {
	return r.setRequestStatus(ctx, supportRequestID, model.ClosedRequestStatus)
}

func (r SupportRequestRepo) GetRequests(ctx context.Context) ([]model.SupportRequest, error) {
	var requests []model.SupportRequest

	rows, err := r.db.Select(ctx, GetRequests)
	if err != nil {
		return nil, err
	}

	if err := pgxscan.ScanAll(requests, rows); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r SupportRequestRepo) GetRequestsWithStatus(ctx context.Context, status model.RequestStatus) ([]model.SupportRequest, error) {
	var requests []model.SupportRequest
	args := pgx.NamedArgs{
		"status": string(status),
	}

	rows, err := r.db.Select(ctx, GetRequests, args)
	if err != nil {
		return nil, err
	}

	if err := pgxscan.ScanAll(requests, rows); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r SupportRequestRepo) setRequestStatus(
	ctx context.Context,
	supportRequestID int,
	status model.RequestStatus,
) error {
	args := pgx.NamedArgs{
		"request_id": supportRequestID,
		"status":     status,
	}
	err := r.db.Update(ctx, SetSupportRequest, args)
	return err
}
