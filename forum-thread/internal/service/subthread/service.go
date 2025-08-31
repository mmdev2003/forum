package subthread

import (
	"context"
	"forum-thread/internal/model"
)

func New(
	subthreadRepo model.ISubthreadRepo,
	messageBroker model.IMessageBroker,
) *ServiceSubthread {
	return &ServiceSubthread{
		subthreadRepo,
		messageBroker,
	}
}

type ServiceSubthread struct {
	subthreadRepo model.ISubthreadRepo
	messageBroker model.IMessageBroker
}

func (s *ServiceSubthread) CreateSubthread(ctx context.Context,
	threadID int,
	threadName,
	subthreadName,
	subthreadDescription string,
) (int, error) {
	subthreadID, err := s.subthreadRepo.CreateSubthread(ctx, threadID, threadName, subthreadName, subthreadDescription)
	if err != nil {
		return 0, err
	}

	return subthreadID, nil
}

func (s *ServiceSubthread) AddViewToSubthread(ctx context.Context,
	subthreadID int,
) error {
	err := s.messageBroker.Publish(ctx, model.AddViewToSubthreadQueue, AddViewToSubthreadPostprocessingBody{
		SubthreadID: subthreadID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceSubthread) AddViewToSubthreadPostprocessing(ctx context.Context,
	subthreadID int,
) error {
	err := s.subthreadRepo.AddViewToSubthread(ctx, subthreadID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceSubthread) SubthreadsByThreadID(ctx context.Context,
	threadID int,
) ([]*model.Subthread, error) {
	subthreads, err := s.subthreadRepo.SubthreadsByThreadID(ctx, threadID)
	if err != nil {
		return nil, err
	}

	return subthreads, nil
}
