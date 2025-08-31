package filter

import (
	"context"
	"errors"
	"fmt"
	"forum-notification/internal/model"

	"github.com/jackc/pgx/v5"
)

type RepoNotificationFilter struct {
	db model.IDatabase
}

func New(db model.IDatabase) *RepoNotificationFilter {
	return &RepoNotificationFilter{
		db: db,
	}
}

func (r *RepoNotificationFilter) IsNotificationEnabled(ctx context.Context, accountID int, notificationType model.NotificationType) (bool, error) {
	args := pgx.NamedArgs{
		"account_id":        accountID,
		"notification_type": string(notificationType),
	}

	var enabled bool
	rows, err := r.db.Select(ctx, CheckIfNotificationTypeIsEnabled, args)
	if err != nil {
		return false, fmt.Errorf("failed to check notification status: %w", err)
	}
	err = rows.Scan(&enabled)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return true, nil
		}
		return false, fmt.Errorf("failed to check notification status: %w", err)
	}

	return enabled, nil
}

func (r *RepoNotificationFilter) GetFilters(ctx context.Context, accountID int) ([]model.NotificationType, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	var rawTypes []string
	rows, err := r.db.Select(ctx, GetFiltersByAccountId, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification enabled_types: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("no notification settings found for account %d", accountID)
	}

	err = rows.Scan(&rawTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse notification enabled_types: %w", err)
	}

	enabledTypes := make([]model.NotificationType, len(rawTypes))
	for i, t := range rawTypes {
		enabledTypes[i] = model.NotificationType(t)
	}

	return enabledTypes, nil
}

func (r *RepoNotificationFilter) UpsertFilters(ctx context.Context, accountID int, enabledTypes []model.NotificationType) error {
	stringSlice := make([]string, len(enabledTypes))
	for i, t := range enabledTypes {
		stringSlice[i] = string(t)
	}

	fmt.Println("Setting filters:", stringSlice)
	args := pgx.NamedArgs{
		"account_id":    accountID,
		"enabled_types": stringSlice,
	}
	err := r.db.Update(ctx, UpsertFilterByAccountId, args)
	if err != nil {
		return fmt.Errorf("failed to upsert notification enabled_types: %w", err)
	}

	return nil
}
