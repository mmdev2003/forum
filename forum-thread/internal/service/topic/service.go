package subthread

import (
	"context"
	"forum-thread/internal/model"
	"forum-thread/pkg/api/forum-status"
)

func New(
	topicRepo model.ITopicRepo,
	accountStatisticRepo model.IAccountStatisticRepo,
	messageBroker model.IMessageBroker,
	notificationClient model.INotificationClient,
	statusClient model.IStatusClient,
) *ServiceTopic {
	return &ServiceTopic{
		topicRepo,
		accountStatisticRepo,
		messageBroker,
		notificationClient,
		statusClient,
	}
}

type ServiceTopic struct {
	topicRepo            model.ITopicRepo
	accountStatisticRepo model.IAccountStatisticRepo
	messageBroker        model.IMessageBroker
	notificationClient   model.INotificationClient
	statusClient         model.IStatusClient
}

func (s *ServiceTopic) CreateTopic(ctx context.Context,
	subthreadID,
	threadID,
	topicOwnerAccountID int,
	threadName,
	subthreadName,
	topicName,
	topicOwnerLogin string,
) (int, error) {
	accountStatuses, err := s.statusClient.StatusByAccountID(ctx, topicOwnerAccountID)
	if err != nil {
		return 0, err
	}
	mergedPermissions := mergeStatusPermissions(accountStatuses)

	if !contains(mergedPermissions.PrivateThreads, threadName) {
		return 0, model.ErrNoPermissionToCreateTopicInThisThread
	}
	if !contains(mergedPermissions.PrivateSubthreads, subthreadName) {
		return 0, model.ErrNoPermissionToCreateTopicInThisSubthread
	}

	todayTopics, err := s.topicRepo.TopicsByAccountIDToday(ctx, topicOwnerAccountID)
	if err != nil {
		return 0, err
	}

	if len(todayTopics) >= mergedPermissions.MaxTopicsPerDay {
		return 0, model.ErrMaxTopicsPerDay
	}

	subthreadTopics, err := s.topicRepo.TopicsBySubthreadIDAndAccountID(ctx, subthreadID, topicOwnerAccountID)
	if err != nil {
		return 0, err
	}

	if len(subthreadTopics) >= mergedPermissions.MaxTopicsInSubthread {
		return 0, model.MaxTopicsInSubthread
	}

	var topicIsAuthor bool
	for _, status := range accountStatuses {
		if status.StatusID == 1 {
			topicIsAuthor = true
		}
	}

	topicID, err := s.topicRepo.CreateTopic(ctx,
		subthreadID,
		threadID,
		topicOwnerAccountID,
		subthreadName,
		threadName,
		topicName,
		topicOwnerLogin,
		topicIsAuthor,
	)
	if err != nil {
		return 0, err
	}

	err = s.messageBroker.Publish(ctx, model.CreateTopicQueue, CreateTopicPostprocessingBody{
		TopicOwnerAccountID: topicOwnerAccountID,
	})
	if err != nil {
		return 0, err
	}

	return topicID, nil
}

func (s *ServiceTopic) CreateTopicPostprocessing(ctx context.Context,
	topicOwnerAccountID int,
) error {
	err := s.accountStatisticRepo.AddCreatedTopicsCount(ctx, topicOwnerAccountID)
	if err != nil {
		return err
	}

	topics, err := s.topicRepo.TopicsByAccountID(ctx, topicOwnerAccountID)
	if err != nil {
		return err
	}

	var countActiveTopics int
	for _, topic := range topics {
		if topic.TopicIsClosed == false && topic.TopicModerationStatus == model.ApprovedTopicStatus {
			countActiveTopics++
		}
	}
	if countActiveTopics >= 50 {
		err = s.statusClient.AssignStatusToAccount(ctx, 1, topicOwnerAccountID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceTopic) AddViewToTopic(ctx context.Context,
	topicID int,
) error {
	err := s.messageBroker.Publish(ctx, model.AddViewToTopicQueue, AddViewToTopicPostprocessingBody{
		TopicID: topicID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceTopic) RejectTopic(ctx context.Context,
	topicID int,
) error {
	err := s.topicRepo.RejectTopic(ctx, topicID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceTopic) ApproveTopic(ctx context.Context,
	topicID int,
) error {
	err := s.topicRepo.ApproveTopic(ctx, topicID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceTopic) AddViewToTopicPostprocessing(ctx context.Context,
	topicID int,
) error {
	err := s.topicRepo.AddViewToTopic(ctx, topicID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceTopic) CloseTopic(ctx context.Context,
	topicOwnerAccountID,
	adminAccountID,
	topicID int,
	topicName,
	adminLogin string,
) error {
	err := s.messageBroker.Publish(ctx, model.CloseTopicQueue, CloseTopicPostprocessingBody{
		TopicOwnerAccountID: topicOwnerAccountID,
		AdminAccountID:      adminAccountID,
		TopicID:             topicID,
		TopicName:           topicName,
		AdminLogin:          adminLogin,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceTopic) CloseTopicPostprocessing(ctx context.Context,
	topicOwnerAccountID,
	adminAccountID,
	topicID int,
	topicName,
	adminLogin string,
) error {
	err := s.topicRepo.CloseTopic(ctx, topicID)
	if err != nil {
		return err
	}

	err = s.notificationClient.TopicClosedNotificationRequest(ctx,
		topicOwnerAccountID,
		adminAccountID,
		topicID,
		topicName,
		adminLogin,
	)

	return nil
}

func (s *ServiceTopic) UpdateTopicAvatar(ctx context.Context,
	topicID int,
	fileName,
	extension string,
	fileBytes []byte,
) error {
	fileUrl, err := s.topicRepo.UploadAvatar(ctx, fileBytes, fileName+extension)
	if err != nil {
		return err
	}

	err = s.topicRepo.UpdateTopicAvatar(ctx, topicID, fileUrl)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceTopic) DownloadTopicAvatar(ctx context.Context,
	fileUrl string,
) ([]byte, error) {
	fileBytes, err := s.topicRepo.DownloadAvatar(ctx, fileUrl)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func (s *ServiceTopic) ChangeTopicPriority(ctx context.Context,
	subthreadID,
	topicID int,
	topicPriority int,
) error {
	err := s.messageBroker.Publish(ctx, model.ChangeTopicPriorityQueue, ChangeTopicPriorityPostprocessingBody{
		SubthreadID:   subthreadID,
		TopicID:       topicID,
		TopicPriority: topicPriority,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceTopic) ChangeTopicPriorityPostprocessing(ctx context.Context,
	subthreadID,
	topicID int,
	topicPriority int,
) error {
	return nil
}

func (s *ServiceTopic) TopicsBySubthreadID(
	ctx context.Context,
	subthreadID int,
) ([]*model.Topic, error) {
	topics, err := s.topicRepo.TopicsBySubthreadID(ctx, subthreadID)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (s *ServiceTopic) TopicsByAccountID(ctx context.Context,
	accountID int,
) ([]*model.Topic, error) {
	topics, err := s.topicRepo.TopicsByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (s *ServiceTopic) TopicsOnModeration(ctx context.Context) ([]*model.Topic, error) {
	topics, err := s.topicRepo.TopicsOnModeration(ctx)
	if err != nil {
		return nil, err
	}
	return topics, nil
}

func mergeStatusPermissions(statuses []*forum_status.Status) model.StatusPermission {
	var merged model.StatusPermission

	for _, status := range statuses {
		perm, exists := model.StatusPermissionMap[status.StatusID]
		if !exists {
			continue
		}

		merged.PrivateThreads = append(merged.PrivateThreads, perm.PrivateThreads...)
		merged.PrivateSubthreads = append(merged.PrivateSubthreads, perm.PrivateSubthreads...)
		merged.PrivateTopics = append(merged.PrivateTopics, perm.PrivateTopics...)

		if perm.MaxTopicsPerDay > merged.MaxTopicsPerDay {
			merged.MaxTopicsPerDay = perm.MaxTopicsPerDay
		}
		if perm.MaxTopicsInSubthread > merged.MaxTopicsInSubthread {
			merged.MaxTopicsInSubthread = perm.MaxTopicsInSubthread
		}
	}

	merged.PrivateThreads = unique(merged.PrivateThreads)
	merged.PrivateSubthreads = unique(merged.PrivateSubthreads)
	merged.PrivateTopics = unique(merged.PrivateTopics)

	return merged
}

func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
