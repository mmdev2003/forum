package subthread

import (
	"context"
	"forum-thread/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func New(
	db model.IDatabase,
	fullTextSearchEngine model.IFullTextSearchEngine,
) *RepoSubthread {
	return &RepoSubthread{
		db:                   db,
		fullTextSearchEngine: fullTextSearchEngine,
	}
}

type RepoSubthread struct {
	db                   model.IDatabase
	fullTextSearchEngine model.IFullTextSearchEngine
}

func (r *RepoSubthread) CreateSubthread(ctx context.Context,
	threadID int,
	threadName,
	subthreadName,
	subthreadDescription string,
) (int, error) {
	args := pgx.NamedArgs{
		"thread_id":             threadID,
		"thread_name":           threadName,
		"subthread_name":        subthreadName,
		"subthread_description": subthreadDescription,
	}
	subthreadID, err := r.db.Insert(ctx, CreateSubthread, args)
	if err != nil {
		return 0, err
	}

	return subthreadID, nil
}

func (r *RepoSubthread) AddViewToSubthread(ctx context.Context,
	subthreadID int,
) error {
	args := pgx.NamedArgs{
		"subthread_id": subthreadID,
	}
	err := r.db.Update(ctx, AddViewToSubthread, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoSubthread) AddMessageCountToSubthread(ctx context.Context,
	subthreadID int,
) error {
	args := pgx.NamedArgs{
		"subthread_id": subthreadID,
	}
	err := r.db.Update(ctx, AddMessageCountToSubthread, args)
	if err != nil {
		return err
	}

	return nil
}
func (r *RepoSubthread) UpdateSubthreadLastMessage(ctx context.Context,
	subthreadID int,
	subthreadLastMessageLogin,
	subthreadLastMessageText string,
) error {
	args := pgx.NamedArgs{
		"subthread_id":                 subthreadID,
		"subthread_last_message_login": subthreadLastMessageLogin,
		"subthread_last_message_text":  subthreadLastMessageText,
	}
	err := r.db.Update(ctx, UpdateSubthreadLastMessage, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepoSubthread) SubthreadsByThreadID(ctx context.Context,
	threadID int,
) ([]*model.Subthread, error) {
	args := pgx.NamedArgs{
		"thread_id": threadID,
	}
	rows, err := r.db.Select(ctx, SubthreadsByThreadID, args)
	if err != nil {
		return nil, err
	}

	var subthreads []*model.Subthread
	err = pgxscan.ScanAll(&subthreads, rows)
	if err != nil {
		return nil, err
	}

	return subthreads, nil
}
