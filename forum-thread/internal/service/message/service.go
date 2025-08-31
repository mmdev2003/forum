package message

import (
	"context"
	"errors"
	"forum-thread/internal/model"
	"forum-thread/pkg/api/forum-status"
	"regexp"
	"time"
)

func New(
	subthreadRepo model.ISubthreadRepo,
	topicRepo model.ITopicRepo,
	messageRepo model.IMessageRepo,
	accountStatisticRepo model.IAccountStatisticRepo,
	supportClient model.ISupportClient,
	statusClient model.IStatusClient,
	notificationClient model.INotificationClient,
	userClient model.IUserClient,
	messageBroker model.IMessageBroker,
) *ServiceMessage {
	return &ServiceMessage{
		subthreadRepo:        subthreadRepo,
		topicRepo:            topicRepo,
		messageRepo:          messageRepo,
		accountStatisticRepo: accountStatisticRepo,
		supportClient:        supportClient,
		statusClient:         statusClient,
		notificationClient:   notificationClient,
		userClient:           userClient,
		messageBroker:        messageBroker,
	}
}

type ServiceMessage struct {
	subthreadRepo        model.ISubthreadRepo
	topicRepo            model.ITopicRepo
	messageRepo          model.IMessageRepo
	accountStatisticRepo model.IAccountStatisticRepo
	supportClient        model.ISupportClient
	statusClient         model.IStatusClient
	notificationClient   model.INotificationClient
	userClient           model.IUserClient
	messageBroker        model.IMessageBroker
}

func (s *ServiceMessage) SendMessageToTopic(ctx context.Context,
	subthreadID,
	topicID,
	replyToMessageID,
	replyMessageOwnerAccountID,
	topicOwnerAccountID,
	senderAccountID int,
	senderLogin,
	threadName,
	subthreadName,
	topicName,
	senderMessageText string,
	uploadedFiles [][]byte,
	filesNames []string,
	filesExtensions []string,
) ([]string, error) {
	accountStatuses, err := s.statusClient.StatusByAccountID(ctx, senderAccountID)
	if err != nil {
		return nil, err
	}
	mergedPermissions := mergeStatusPermissions(accountStatuses)

	if !contains(mergedPermissions.PrivateThreads, threadName) {
		return nil, errors.New("user is not allowed to create topics in thread: " + threadName)
	}
	if !contains(mergedPermissions.PrivateSubthreads, subthreadName) {
		return nil, errors.New("user is not allowed to create topics in subthread: " + subthreadName)
	}

	var filesURLs []string
	var filesSizes []int
	for i, _ := range uploadedFiles {
		fileURL, err := s.messageRepo.UploadFile(ctx, uploadedFiles[i], filesNames[i])
		if err != nil {
			return nil, err
		}
		filesURLs = append(filesURLs, fileURL)
		filesSizes = append(filesSizes, len(uploadedFiles[i]))
	}

	err = s.messageBroker.Publish(ctx, model.SendMessageToTopicQueue, SendMessageToTopicPostprocessingBody{
		SubthreadID:                subthreadID,
		TopicID:                    topicID,
		ReplyToMessageID:           replyToMessageID,
		ReplyMessageOwnerAccountID: replyMessageOwnerAccountID,
		TopicOwnerAccountID:        topicOwnerAccountID,
		SenderAccountID:            senderAccountID,
		SenderLogin:                senderLogin,
		SenderMessageText:          senderMessageText,
		TopicName:                  topicName,
		FilesURLs:                  filesURLs,
		FilesNames:                 filesNames,
		FilesExtensions:            filesExtensions,
		FilesSizes:                 filesSizes,
	})
	if err != nil {
		return nil, err
	}

	return filesURLs, nil
}

func (s *ServiceMessage) SendMessageToTopicPostprocessing(ctx context.Context,
	subthreadID,
	topicID,
	replyToMessageID,
	replyMessageOwnerAccountID,
	topicOwnerAccountID,
	senderAccountID int,
	senderLogin,
	topicName,
	senderMessageText string,
	filesURLs []string,
	filesNames []string,
	filesExtensions []string,
	filesSizes []int,
) (int, error) {
	ctx, err := s.messageRepo.CtxWithTx(ctx)
	if err != nil {
		return 0, err
	}

	senderMessageID, err := s.messageRepo.CreateMessage(
		ctx,
		topicID,
		replyToMessageID,
		senderAccountID,
		senderLogin,
		senderMessageText,
	)
	if err != nil {
		s.messageRepo.RollbackTx(ctx)
		return 0, err
	}
	for i, _ := range filesURLs {
		err := s.messageRepo.AddFileToMessage(ctx,
			senderMessageID,
			filesSizes[i],
			filesURLs[i],
			filesNames[i],
			filesExtensions[i],
		)
		if err != nil {
			return 0, err
		}
	}

	err = s.accountStatisticRepo.AddSentMessagesToTopicsCount(ctx, senderAccountID)
	if err != nil {
		return 0, err
	}

	err = s.subthreadRepo.AddMessageCountToSubthread(ctx, topicID)
	if err != nil {
		s.messageRepo.RollbackTx(ctx)
		return 0, err
	}

	err = s.subthreadRepo.UpdateSubthreadLastMessage(ctx, subthreadID, senderLogin, senderMessageText)
	if err != nil {
		s.messageRepo.RollbackTx(ctx)
		return 0, err
	}

	err = s.topicRepo.AddMessageCountToTopic(ctx, topicID)
	if err != nil {
		s.messageRepo.RollbackTx(ctx)
		return 0, err
	}

	err = s.topicRepo.UpdateTopicLastMessage(ctx, topicID, senderLogin, senderMessageText)
	if err != nil {
		s.messageRepo.RollbackTx(ctx)
		return 0, err
	}

	err = s.messageRepo.CreateMessageSearchIndex(ctx,
		topicID,
		senderMessageID,
		senderAccountID,
		senderLogin,
		senderMessageText,
	)
	if err != nil {
		s.messageRepo.RollbackTx(ctx)
		return 0, err
	}

	messages, err := s.messageRepo.MessagesByAccountID(ctx, senderAccountID)
	if err != nil {
		return 0, err
	}

	var lastMonthMessagesCount int
	for _, message := range messages {
		oneMonthAgo := time.Now().AddDate(0, -1, 0)
		if message.CreatedAt.After(oneMonthAgo) {
			lastMonthMessagesCount++
		}
	}

	if lastMonthMessagesCount >= 100 {
		err = s.statusClient.AssignStatusToAccount(ctx, 2, senderAccountID)
		if err != nil {
			return 0, err
		}
	} else {
		err = s.statusClient.RevokeStatusFromAccount(ctx, 2, senderAccountID)
		if err != nil {
			return 0, err
		}
	}

	err = s.messageRepo.CommitTx(ctx)
	if err != nil {
		s.messageRepo.RollbackTx(ctx)
		return 0, err
	}

	err = s.notificationClient.NewMessageFromTopicNotification(ctx,
		topicOwnerAccountID,
		senderMessageID,
		senderAccountID,
		topicID,
		senderMessageText,
		topicName,
		senderLogin,
	)
	if err != nil {
		return 0, err
	}

	mentions := s.parseLoginFromTopicMessage(senderMessageText)
	for _, mentionLogin := range mentions {
		user, err := s.userClient.UserByLogin(ctx, mentionLogin)
		if err != nil {
			return 0, err
		}

		mentionedAccountID := user.User[0].AccountID
		err = s.notificationClient.NewMentionFromTopicNotification(ctx,
			mentionedAccountID,
			senderMessageID,
			senderAccountID,
			senderMessageText,
			topicName,
			senderLogin,
		)
		if err != nil {
			return 0, err
		}
	}
	if replyToMessageID != 0 {
		err = s.notificationClient.NewMessageReplyFromTopicNotification(ctx,
			replyMessageOwnerAccountID,
			senderMessageID,
			senderAccountID,
			topicID,
			senderMessageText,
			topicName,
			senderLogin,
		)
		if err != nil {
			return 0, err
		}
	}

	return senderMessageID, nil
}

func (s *ServiceMessage) LikeMessage(ctx context.Context,
	topicID,
	messageOwnerAccountID,
	likerAccountID,
	likeMessageID,
	likeTypeID int,
	likerLogin,
	topicName,
	likeMessageText string,
) error {
	err := s.messageBroker.Publish(ctx, model.LikeMessageQueue, LikeMessagePostprocessingBody{
		TopicID:               topicID,
		MessageOwnerAccountID: messageOwnerAccountID,
		LikerAccountID:        likerAccountID,
		LikeMessageID:         likeMessageID,
		LikeTypeID:            likeTypeID,
		LikerLogin:            likerLogin,
		TopicName:             topicName,
		LikeMessageText:       likeMessageText,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceMessage) LikeMessagePostprocessing(ctx context.Context,
	topicID,
	messageOwnerAccountID,
	likerAccountID,
	likeMessageID,
	likeTypeID int,
	likerLogin,
	topicName,
	likeMessageText string,
) error {
	ctx, err := s.messageRepo.CtxWithTx(ctx)
	if err != nil {
		return err
	}

	_, err = s.messageRepo.CreateMessageLike(ctx, topicID, likeMessageID, likeTypeID, likerAccountID)
	if err != nil {
		s.messageRepo.RollbackTx(ctx)
		return err
	}

	err = s.messageRepo.IncrementLikeCountToMessage(ctx, likeMessageID)
	if err != nil {
		s.messageRepo.RollbackTx(ctx)
		return err
	}

	err = s.messageRepo.CommitTx(ctx)
	if err != nil {
		s.messageRepo.RollbackTx(ctx)
		return err
	}

	err = s.notificationClient.NewLikeMessageFromTopicNotification(ctx,
		messageOwnerAccountID,
		likeMessageID,
		likerAccountID,
		topicID,
		likeMessageText,
		topicName,
		likerLogin,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceMessage) ReportMessage(ctx context.Context,
	messageID,
	accountID int,
	reportText string,
) error {
	err := s.messageBroker.Publish(ctx, model.ReportMessageQueue, ReportMessagePostprocessingBody{
		MessageID:  messageID,
		ReportText: reportText,
		AccountID:  accountID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceMessage) ReportMessagePostprocessing(ctx context.Context,
	messageID,
	accountID int,
	reportText string,
) error {
	ctx, err := s.messageRepo.CtxWithTx(ctx)
	if err != nil {
		return err
	}

	err = s.messageRepo.AddReportCountToMessage(ctx, messageID)
	if err != nil {
		s.messageRepo.RollbackTx(ctx)
		return err
	}

	err = s.messageRepo.CommitTx(ctx)
	if err != nil {
		s.messageRepo.RollbackTx(ctx)
		return err
	}

	err = s.supportClient.ReportMessage(ctx, accountID, messageID, reportText)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceMessage) MessagesByText(ctx context.Context,
	text string,
) ([]*model.MessageSearch, error) {
	messages, err := s.messageRepo.MessagesByText(ctx, text)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *ServiceMessage) MessagesByTopicID(ctx context.Context,
	accountID int,
	topicID int,
) ([]*model.Message, []*model.Like, []*model.File, error) {
	messages, err := s.messageRepo.MessagesByTopicID(ctx, topicID)
	if err != nil {
		return nil, nil, nil, err
	}

	if accountID == 0 {
		return messages, nil, nil, err
	}

	accountLikesOnMessage, err := s.messageRepo.LikesByTopicIDAndAccountID(ctx, topicID, accountID)
	if err != nil {
		return messages, nil, nil, err
	}

	var files []*model.File
	for _, message := range messages {
		filesOnMessage, err := s.messageRepo.FilesByMessageID(ctx, message.ID)
		if err != nil {
			return messages, accountLikesOnMessage, nil, err
		}
		files = append(files, filesOnMessage...)
	}

	return messages, accountLikesOnMessage, files, nil
}

func (s *ServiceMessage) EditMessage(ctx context.Context,
	messageID int,
	messageText string,
) error {
	err := s.messageRepo.EditMessage(ctx, messageID, messageText)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceMessage) MessagesByAccountID(ctx context.Context,
	accountID int,
) ([]*model.Message, []*model.File, error) {
	messages, err := s.messageRepo.MessagesByAccountID(ctx, accountID)
	if err != nil {
		return nil, nil, err
	}

	var files []*model.File
	for _, message := range messages {
		filesOnMessage, err := s.messageRepo.FilesByMessageID(ctx, message.ID)
		if err != nil {
			return messages, nil, err
		}
		files = append(files, filesOnMessage...)
	}

	return messages, files, nil
}

func (s *ServiceMessage) UnlikeMessage(ctx context.Context,
	likeMessageID,
	likerAccountID int,
) error {
	err := s.messageRepo.DeleteMessageLike(ctx, likeMessageID, likerAccountID)
	if err != nil {
		return err
	}

	err = s.messageRepo.DecrementLikeCountToMessage(ctx, likeMessageID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceMessage) DownloadFile(ctx context.Context,
	fileURL string,
) ([]byte, error) {
	fileBytes, err := s.messageRepo.DownloadFile(ctx, fileURL)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func (s *ServiceMessage) parseLoginFromTopicMessage(text string) []string {
	re := regexp.MustCompile(`@([\p{L}\d._-]+)`)
	seen := make(map[string]struct{})
	var mentions []string

	for _, match := range re.FindAllStringSubmatch(text, -1) {
		mention := match[1]
		if _, exists := seen[mention]; !exists {
			seen[mention] = struct{}{}
			mentions = append(mentions, mention)
		}
	}

	return mentions
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
