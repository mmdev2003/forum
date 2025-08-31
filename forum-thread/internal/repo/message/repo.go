package message

import (
	"context"
	"encoding/json"
	"forum-thread/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func New(
	db model.IDatabase,
	fullTextSearchEngine model.IFullTextSearchEngine,
	weedFS model.IWeedFS,
) *RepoMessage {
	return &RepoMessage{
		db:                   db,
		fullTextSearchEngine: fullTextSearchEngine,
		weedFS:               weedFS,
	}
}

type RepoMessage struct {
	db                   model.IDatabase
	fullTextSearchEngine model.IFullTextSearchEngine
	weedFS               model.IWeedFS
}

func (r *RepoMessage) CreateMessage(ctx context.Context,
	topicID,
	replyToMessageID,
	accountID int,
	login,
	text string,
) (int, error) {
	args := pgx.NamedArgs{
		"topic_id":                 topicID,
		"message_reply_to_id":      replyToMessageID,
		"message_owner_account_id": accountID,
		"message_owner_login":      login,
		"message_text":             text,
	}
	messageID, err := r.db.Insert(ctx, CreateMessage, args)
	if err != nil {
		return 0, err
	}

	return messageID, nil
}

func (r *RepoMessage) UploadFile(ctx context.Context,
	file []byte,
	name string,
) (string, error) {
	fileID, err := r.weedFS.Upload(file, name, int64(len(file)), "")
	if err != nil {
		return "", err
	}
	return fileID, nil
}

func (r *RepoMessage) CreateMessageSearchIndex(ctx context.Context,
	topicID,
	messageID,
	accountID int,
	login,
	text string,
) error {
	document := model.MessageSearch{
		ID:        messageID,
		TopicID:   topicID,
		AccountID: accountID,
		Login:     login,
		Text:      text,
	}
	err := r.fullTextSearchEngine.AddDocuments(model.MessageFullTextSearchIndex, []any{document})
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoMessage) AddFileToMessage(ctx context.Context,
	messageID,
	size int,
	url,
	name,
	extension string,
) error {
	args := pgx.NamedArgs{
		"message_id": messageID,
		"size":       size,
		"url":        url,
		"name":       name,
		"extension":  extension,
	}
	_, err := r.db.Insert(ctx, AddFileToMessage, args)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoMessage) CreateMessageLike(ctx context.Context,
	topicID,
	messageID,
	likeTypeID,
	accountID int,
) (int, error) {
	args := pgx.NamedArgs{
		"topic_id":     topicID,
		"message_id":   messageID,
		"like_type_id": likeTypeID,
		"account_id":   accountID,
	}
	likeID, err := r.db.Insert(ctx, CreateMessageLike, args)
	if err != nil {
		return 0, err
	}

	return likeID, nil
}

func (r *RepoMessage) DeleteMessageLike(ctx context.Context,
	likeMessageID,
	likerAccountID int,
) error {
	args := pgx.NamedArgs{
		"message_id": likeMessageID,
		"account_id": likerAccountID,
	}
	err := r.db.Delete(ctx, DeleteMessageLike, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoMessage) IncrementLikeCountToMessage(ctx context.Context,
	messageID int,
) error {
	args := pgx.NamedArgs{
		"message_id": messageID,
	}
	err := r.db.Update(ctx, IncrementLikeCountToMessage, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoMessage) DecrementLikeCountToMessage(ctx context.Context,
	messageID int,
) error {
	args := pgx.NamedArgs{
		"message_id": messageID,
	}
	err := r.db.Update(ctx, DecrementLikeCountToMessage, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoMessage) AddReplyCountToMessage(ctx context.Context,
	messageID int,
) error {
	args := pgx.NamedArgs{
		"message_id": messageID,
	}
	err := r.db.Update(ctx, AddReplyCountToMessage, args)
	if err != nil {
		return err
	}

	return nil
}
func (r *RepoMessage) AddReportCountToMessage(ctx context.Context,
	messageID int,
) error {
	args := pgx.NamedArgs{
		"message_id": messageID,
	}
	err := r.db.Update(ctx, AddReportCountToMessage, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoMessage) EditMessage(ctx context.Context,
	messageID int,
	messageText string,
) error {
	args := pgx.NamedArgs{
		"message_id":   messageID,
		"message_text": messageText,
	}
	err := r.db.Update(ctx, EditMessage, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoMessage) MessagesByText(ctx context.Context,
	text string,
) ([]*model.MessageSearch, error) {
	result, err := r.fullTextSearchEngine.SimpleSearch(model.MessageFullTextSearchIndex, text)
	if err != nil {
		return nil, err
	}

	var messages []*model.MessageSearch
	err = json.Unmarshal(result, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *RepoMessage) MessagesByTopicID(ctx context.Context,
	topicID int,
) ([]*model.Message, error) {
	args := pgx.NamedArgs{
		"topic_id": topicID,
	}
	rows, err := r.db.Select(ctx, MessagesByTopicID, args)
	if err != nil {
		return nil, err
	}

	var messages []*model.Message
	err = pgxscan.ScanAll(&messages, rows)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *RepoMessage) MessagesByAccountID(ctx context.Context,
	accountID int,
) ([]*model.Message, error) {
	args := pgx.NamedArgs{
		"account_id": accountID,
	}
	rows, err := r.db.Select(ctx, MessageByAccountID, args)
	if err != nil {
		return nil, err
	}

	var messages []*model.Message
	err = pgxscan.ScanAll(&messages, rows)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *RepoMessage) LikesByTopicIDAndAccountID(ctx context.Context,
	topicID,
	accountID int,
) ([]*model.Like, error) {
	args := pgx.NamedArgs{
		"topic_id":   topicID,
		"account_id": accountID,
	}
	rows, err := r.db.Select(ctx, LikesByTopicIDAndAccountID, args)
	if err != nil {
		return nil, err
	}

	var likes []*model.Like
	err = pgxscan.ScanAll(&likes, rows)
	if err != nil {
		return nil, err
	}

	return likes, nil
}

func (r *RepoMessage) FilesByMessageID(ctx context.Context,
	messageID int,
) ([]*model.File, error) {
	args := pgx.NamedArgs{
		"message_id": messageID,
	}
	rows, err := r.db.Select(ctx, FilesByMessageID, args)
	if err != nil {
		return nil, err
	}

	var files []*model.File
	err = pgxscan.ScanAll(&files, rows)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (r *RepoMessage) DownloadFile(ctx context.Context,
	fileURL string,
) ([]byte, error) {
	file, err := r.weedFS.Download(fileURL)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (r *RepoMessage) CtxWithTx(ctx context.Context) (context.Context, error) {
	return r.db.CtxWithTx(ctx)
}
func (r *RepoMessage) CommitTx(ctx context.Context) error {
	return r.db.CommitTx(ctx)
}
func (r *RepoMessage) RollbackTx(ctx context.Context) {
	r.db.RollbackTx(ctx)
}
