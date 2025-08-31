package thread

import (
	"context"
	"forum-thread/internal/model"
)

func New(
	threadRepo model.IThreadRepo,
	messageBroker model.IMessageBroker,
) *ServiceThread {
	return &ServiceThread{
		threadRepo:    threadRepo,
		messageBroker: messageBroker,
	}
}

type ServiceThread struct {
	threadRepo    model.IThreadRepo
	messageBroker model.IMessageBroker
}

func (s *ServiceThread) CreateThread(ctx context.Context,
	threadName,
	threadDescription,
	threadColor string,
	allowedStatuses []string,
) (int, error) {
	threadID, err := s.threadRepo.CreateThread(ctx, threadName, threadDescription, threadColor, allowedStatuses)
	if err != nil {
		return 0, err
	}

	return threadID, nil
}

func (s *ServiceThread) AllThreads(
	ctx context.Context,
) ([]*model.Thread, error) {
	threads, err := s.threadRepo.AllThreads(ctx)
	if err != nil {
		return nil, err
	}

	return threads, nil
}
