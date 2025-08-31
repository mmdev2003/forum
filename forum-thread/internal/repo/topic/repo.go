package topic

import (
	"context"
	"forum-thread/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func New(
	db model.IDatabase,
	fullTextSearchEngine model.IFullTextSearchEngine,
	weedFS model.IWeedFS,
) *RepoTopic {
	return &RepoTopic{
		db:                   db,
		fullTextSearchEngine: fullTextSearchEngine,
		weedFS:               weedFS,
	}
}

type RepoTopic struct {
	db                   model.IDatabase
	fullTextSearchEngine model.IFullTextSearchEngine
	weedFS               model.IWeedFS
}

func (r *RepoTopic) CreateTopic(ctx context.Context,
	subthreadID,
	threadID,
	topicOwnerAccountID int,
	subthreadName,
	threadName,
	topicName,
	topicOwnerLogin string,
	topicIsAuthor bool,
) (int, error) {
	args := pgx.NamedArgs{
		"subthread_id":            subthreadID,
		"thread_id":               threadID,
		"topic_owner_account_id":  topicOwnerAccountID,
		"subthread_name":          subthreadName,
		"thread_name":             threadName,
		"topic_name":              topicName,
		"topic_owner_login":       topicOwnerLogin,
		"topic_moderation_status": model.OnModerationTopicStatus,
		"topic_is_author":         topicIsAuthor,
	}
	topicID, err := r.db.Insert(ctx, CreateTopic, args)
	if err != nil {
		return 0, err
	}

	return topicID, nil
}

func (r *RepoTopic) UpdateTopicLastMessage(ctx context.Context,
	topicID int,
	topicLastMessageLogin,
	topicLastMessageText string,
) error {
	args := pgx.NamedArgs{
		"topic_id":                 topicID,
		"topic_last_message_login": topicLastMessageLogin,
		"topic_last_message_text":  topicLastMessageText,
	}
	err := r.db.Update(ctx, UpdateTopicLastMessage, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoTopic) RejectTopic(ctx context.Context,
	topicID int,
) error {
	args := pgx.NamedArgs{
		"topic_id":                topicID,
		"topic_moderation_status": model.RejectedTopicStatus,
	}
	err := r.db.Update(ctx, RejectTopic, args)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoTopic) ApproveTopic(ctx context.Context,
	topicID int,
) error {
	args := pgx.NamedArgs{
		"topic_id":                topicID,
		"topic_moderation_status": model.ApprovedTopicStatus,
	}
	err := r.db.Update(ctx, ApproveTopic, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoTopic) AddMessageCountToTopic(ctx context.Context,
	topicID int,
) error {
	args := pgx.NamedArgs{
		"topic_id": topicID,
	}
	err := r.db.Update(ctx, AddMessageCountToTopic, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoTopic) AddViewToTopic(ctx context.Context,
	topicID int,
) error {
	args := pgx.NamedArgs{
		"topic_id": topicID,
	}
	err := r.db.Update(ctx, AddViewToTopic, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoTopic) CloseTopic(ctx context.Context,
	topicID int,
) error {
	args := pgx.NamedArgs{
		"topic_id": topicID,
	}
	err := r.db.Update(ctx, CloseTopic, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoTopic) ChangeTopicPriority(ctx context.Context,
	topicID,
	topicPriority int,
) error {
	args := pgx.NamedArgs{
		"topic_id":       topicID,
		"topic_priority": topicPriority,
	}
	err := r.db.Update(ctx, ChangeTopicPriority, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoTopic) TopicsBySubthreadID(ctx context.Context,
	subthreadID int,
) ([]*model.Topic, error) {
	args := pgx.NamedArgs{
		"subthread_id":            subthreadID,
		"topic_moderation_status": model.ApprovedTopicStatus,
	}
	rows, err := r.db.Select(ctx, TopicsBySubthreadID, args)
	if err != nil {
		return nil, err
	}

	var topics []*model.Topic
	err = pgxscan.ScanAll(&topics, rows)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *RepoTopic) TopicsBySubthreadIDAndAccountID(ctx context.Context,
	subthreadID,
	accountID int,
) ([]*model.Topic, error) {
	args := pgx.NamedArgs{
		"subthread_id":            subthreadID,
		"account_id":              accountID,
		"topic_moderation_status": model.ApprovedTopicStatus,
	}
	rows, err := r.db.Select(ctx, TopicsBySubthreadIDAndAccountID, args)
	if err != nil {
		return nil, err
	}

	var topics []*model.Topic
	err = pgxscan.ScanAll(&topics, rows)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *RepoTopic) TopicsByAccountID(ctx context.Context,
	accountID int,
) ([]*model.Topic, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	rows, err := r.db.Select(ctx, TopicsByAccountID, args)
	if err != nil {
		return nil, err
	}

	var topics []*model.Topic
	err = pgxscan.ScanAll(&topics, rows)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *RepoTopic) TopicsOnModeration(ctx context.Context) ([]*model.Topic, error) {
	args := pgx.NamedArgs{
		"topic_moderation_status": model.OnModerationTopicStatus,
	}
	rows, err := r.db.Select(ctx, TopicsOnModeration, args)
	if err != nil {
		return nil, err
	}

	var topics []*model.Topic
	err = pgxscan.ScanAll(&topics, rows)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *RepoTopic) TopicsByAccountIDToday(ctx context.Context,
	accountID int,
) ([]*model.Topic, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	rows, err := r.db.Select(ctx, TopicsByAccountIDToday, args)
	if err != nil {
		return nil, err
	}

	var topics []*model.Topic
	err = pgxscan.ScanAll(&topics, rows)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r *RepoTopic) UpdateTopicAvatar(ctx context.Context,
	topicID int,
	topicAvatarURL string,
) error {
	args := pgx.NamedArgs{
		"topic_id":         topicID,
		"topic_avatar_url": topicAvatarURL,
	}
	err := r.db.Update(ctx, UpdateTopicAvatar, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoTopic) UploadAvatar(ctx context.Context,
	file []byte,
	name string,
) (string, error) {
	fileID, err := r.weedFS.Upload(file, name, int64(len(file)), "")
	if err != nil {
		return "", err
	}
	return fileID, nil
}

func (r *RepoTopic) DownloadAvatar(ctx context.Context,
	fileURL string,
) ([]byte, error) {
	file, err := r.weedFS.Download(fileURL)
	if err != nil {
		return nil, err
	}
	return file, nil
}
