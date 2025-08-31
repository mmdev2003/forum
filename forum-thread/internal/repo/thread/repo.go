package thread

import (
	"context"
	"forum-thread/internal/model"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func New(
	db model.IDatabase,
	fullTextSearchEngine model.IFullTextSearchEngine,
) *RepoThread {
	return &RepoThread{
		db:                   db,
		fullTextSearchEngine: fullTextSearchEngine,
	}
}

type RepoThread struct {
	db                   model.IDatabase
	fullTextSearchEngine model.IFullTextSearchEngine
}

func (r *RepoThread) CreateThread(ctx context.Context,
	threadName,
	threadDescription,
	threadColor string,
	allowedStatuses []string,
) (int, error) {
	args := pgx.NamedArgs{
		"thread_name":        threadName,
		"thread_description": threadDescription,
		"allowed_statuses":   allowedStatuses,
		"thread_color":       threadColor,
	}
	threadID, err := r.db.Insert(ctx, CreateThread, args)
	if err != nil {
		return 0, err
	}

	return threadID, nil
}

func (r *RepoThread) AllThreads(
	ctx context.Context,
) ([]*model.Thread, error) {
	rows, err := r.db.Select(ctx, AllThreads, pgx.NamedArgs{})
	if err != nil {
		return nil, err
	}

	var threads []*model.Thread
	err = pgxscan.ScanAll(&threads, rows)
	if err != nil {
		return nil, err
	}

	return threads, nil
}
