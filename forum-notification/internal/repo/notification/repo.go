package notification

import (
	"context"
	"database/sql"
	"fmt"
	"forum-notification/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
)

func New(db model.IDatabase) *RepoNotification {
	return &RepoNotification{
		db: db,
	}
}

type RepoNotification struct {
	db model.IDatabase
}

func (r *RepoNotification) CreateMessageFromTopicNotification(ctx context.Context,
	accountID,
	messageID,
	replierAccountID,
	topicID int,
	messageText,
	topicName,
	replierLogin string,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id":         accountID,
		"message_id":         messageID,
		"replier_account_id": replierAccountID,
		"topic_id":           topicID,
		"message_text":       messageText,
		"topic_name":         topicName,
		"replier_login":      replierLogin,
	}
	notificationID, err := r.db.Insert(ctx, CreateMessageFromTopicNotification, args)
	if err != nil {
		return 0, err
	}
	return notificationID, err
}

func (r *RepoNotification) CreateMessageReplyFromTopicNotification(ctx context.Context,
	accountID,
	messageID,
	replierAccountID,
	topicID int,
	messageText,
	topicName,
	replierLogin string,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id":         accountID,
		"message_id":         messageID,
		"replier_account_id": replierAccountID,
		"topic_id":           topicID,
		"message_text":       messageText,
		"topic_name":         topicName,
		"replier_login":      replierLogin,
	}
	notificationID, err := r.db.Insert(ctx, CreateMessageReplyFromTopicNotification, args)
	if err != nil {
		return 0, err
	}
	return notificationID, err
}

func (r *RepoNotification) CreateLikeMessageFromTopicNotification(ctx context.Context,
	accountID,
	messageID,
	likerAccountID,
	topicID int,
	messageText,
	topicName,
	likerLogin string,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id":       accountID,
		"message_id":       messageID,
		"liker_account_id": likerAccountID,
		"topic_id":         topicID,
		"message_text":     messageText,
		"topic_name":       topicName,
		"liker_login":      likerLogin,
	}
	notificationID, err := r.db.Insert(ctx, CreateLikeMessageFromTopicNotification, args)
	if err != nil {
		return 0, err
	}
	return notificationID, err
}

func (r *RepoNotification) CreateTopicClosedNotification(ctx context.Context,
	accountID,
	adminAccountID,
	topicID int,
	topicName,
	adminLogin string,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id":       accountID,
		"admin_account_id": adminAccountID,
		"topic_id":         topicID,
		"topic_name":       topicName,
		"admin_login":      adminLogin,
	}
	notificationID, err := r.db.Insert(ctx, CreateTopicClosedNotification, args)
	if err != nil {
		return 0, err
	}
	return notificationID, err
}

func (r *RepoNotification) CreateResponseToSupportRequestNotification(ctx context.Context,
	accountID,
	supportRequestID int,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id":         accountID,
		"support_request_id": supportRequestID,
	}
	notificationID, err := r.db.Insert(ctx, CreateResponseToSupportRequestNotification, args)
	if err != nil {
		return 0, err
	}
	return notificationID, err
}

func (r *RepoNotification) CreateStatusReceivedNotification(ctx context.Context,
	accountID int,
	statusName string,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id":  accountID,
		"status_name": statusName,
	}
	notificationID, err := r.db.Insert(ctx, CreateStatusReceivedNotification, args)
	if err != nil {
		return 0, err
	}
	return notificationID, err
}

func (r *RepoNotification) CreateFrameReceivedNotification(ctx context.Context,
	accountID int,
	frameName string,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
		"frame_name": frameName,
	}
	notificationID, err := r.db.Insert(ctx, CreateFrameReceivedNotification, args)
	if err != nil {
		return 0, err
	}
	return notificationID, err
}

func (r *RepoNotification) CreateMessageFromDialogNotification(ctx context.Context,
	accountID,
	messageID,
	dialogID,
	senderAccountID int,
	messageText,
	senderLogin string,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id":        accountID,
		"message_id":        messageID,
		"dialog_id":         dialogID,
		"sender_account_id": senderAccountID,
		"message_text":      messageText,
		"sender_login":      senderLogin,
	}
	notificationID, err := r.db.Insert(ctx, CreateMessageFromDialogNotification, args)
	if err != nil {
		return 0, err
	}
	return notificationID, err
}

func (r *RepoNotification) CreateMentionFromTopicNotification(ctx context.Context,
	accountID,
	messageID,
	mentionAccountID int,
	messageText,
	topicName,
	mentionLogin string,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id":         accountID,
		"message_id":         messageID,
		"mention_account_id": mentionAccountID,
		"message_text":       messageText,
		"topic_name":         topicName,
		"mention_login":      mentionLogin,
	}
	notificationID, err := r.db.Insert(ctx, CreateMentionFromTopicNotification, args)
	if err != nil {
		return 0, err
	}
	return notificationID, err
}

func (r *RepoNotification) CreateWarningFromAdminNotification(ctx context.Context,
	accountID,
	adminAccountID int,
	warningText,
	adminLogin string,
) (int, error) {
	args := pgx.NamedArgs{
		"account_id":       accountID,
		"admin_account_id": adminAccountID,
		"message_text":     warningText,
		"admin_login":      adminLogin,
	}
	notificationID, err := r.db.Insert(ctx, CreateWarningFromAdminNotification, args)
	if err != nil {
		return 0, err
	}
	return notificationID, err
}

func (r *RepoNotification) GetNotificationsByAccountID(ctx context.Context, id int) ([]interface{}, error) {
	args := pgx.NamedArgs{"account_id": id}
	rows, err := r.db.Select(ctx, GetNotificationsByAccountID, args)
	if err != nil {
		return nil, err
	}

	var notifications []interface{}

	for rows.Next() {
		notification, err := ParseNotificationRow(rows)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

func ParseNotificationRow(rows pgx.Rows) (interface{}, error) {
	var (
		notificationID, accountID                                            int
		notificationType                                                     model.NotificationType
		isRead                                                               bool
		createdAt, updatedAt                                                 time.Time
		messageID, replierAccountID, topicID, likerAccountID, adminAccountID sql.NullInt64
		supportRequestID, supportRequestIDClosed, dialogID, senderAccountID  sql.NullInt64
		mentionAccountID, warningAdminAccountID                              sql.NullInt64
		messageText, topicName, replierLogin, likerLogin, adminLogin         sql.NullString
		statusName, frameName, senderLogin, mentionLogin, warningAdminLogin  sql.NullString
	)
	err := rows.Scan(
		&notificationID, &notificationType, &isRead, &createdAt, &updatedAt,
		&accountID,
		&messageID, &replierAccountID, &topicID,
		&messageText, &topicName, &replierLogin,
		&likerAccountID, &likerLogin,
		&adminAccountID, &adminLogin,
		&supportRequestID, &supportRequestIDClosed,
		&statusName, &frameName,
		&dialogID, &senderAccountID, &senderLogin,
		&mentionAccountID, &mentionLogin,
		&warningAdminAccountID, &warningAdminLogin,
	)
	if err != nil {
		return nil, err
	}

	notification := model.Notification{
		ID:        notificationID,
		Type:      notificationType,
		IsRead:    isRead,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	switch notification.Type {
	case model.MessageFromTopicType:
		return model.NewMessageFromTopicNotification{
			Notification:      notification,
			SenderMessageID:   int(messageID.Int64),
			SenderAccountID:   int(replierAccountID.Int64),
			TopicID:           int(topicID.Int64),
			SenderMessageText: messageText.String,
			TopicName:         topicName.String,
			SenderLogin:       replierLogin.String,
		}, nil
	case model.MessageReplyFromTopicType:
		return model.NewMessageReplyFromTopicNotification{
			Notification:      notification,
			SenderMessageID:   int(messageID.Int64),
			SenderAccountID:   int(replierAccountID.Int64),
			TopicID:           int(topicID.Int64),
			SenderMessageText: messageText.String,
			TopicName:         topicName.String,
			SenderLogin:       replierLogin.String,
		}, nil
	case model.LikeMessageFromTopicType:
		return model.NewLikeMessageFromTopicNotification{
			Notification:      notification,
			SenderMessageID:   int(messageID.Int64),
			LikerAccountID:    int(likerAccountID.Int64),
			TopicID:           int(topicID.Int64),
			SenderMessageText: messageText.String,
			TopicName:         topicName.String,
			LikerLogin:        likerLogin.String,
		}, nil
	case model.TopicClosedType:
		return model.TopicClosedNotification{
			Notification:   notification,
			AdminAccountID: int(adminAccountID.Int64),
			TopicID:        int(topicID.Int64),
			TopicName:      topicName.String,
			AdminLogin:     adminLogin.String,
		}, nil
	case model.ResponseToSupportRequestType:
		return model.ResponseToSupportRequestNotification{
			Notification:     notification,
			SupportRequestID: int(supportRequestID.Int64),
		}, nil
	case model.StatusReceivedType:
		return model.StatusReceivedNotification{
			Notification: notification,
			StatusName:   statusName.String,
		}, nil
	case model.FrameReceivedType:
		return model.FrameReceivedNotification{
			Notification: notification,
			FrameName:    frameName.String,
		}, nil
	case model.MentionFromTopicType:
		return model.NewMentionFromTopicNotification{
			Notification:      notification,
			MessageID:         int(messageID.Int64),
			SenderAccountID:   int(mentionAccountID.Int64),
			SenderMessageText: messageText.String,
			TopicName:         topicName.String,
			SenderLogin:       mentionLogin.String,
		}, nil
	case model.WarningFromAdminType:
		return model.NewWarningFromAdminNotification{
			Notification:   notification,
			AdminAccountID: int(warningAdminAccountID.Int64),
			WarningText:    messageText.String,
			AdminLogin:     warningAdminLogin.String,
		}, nil
	case model.MessageFromDialogType:
		return model.NewMessageFromDialogNotification{
			Notification:    notification,
			MessageID:       int(messageID.Int64),
			DialogID:        int(dialogID.Int64),
			SenderAccountID: int(senderAccountID.Int64),
			MessageText:     messageText.String,
			SenderLogin:     senderLogin.String,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported notification type: %s", notification.Type)
	}
}
