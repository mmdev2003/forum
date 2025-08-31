package model

import (
	"context"
	"forum-thread/pkg/api/forum-authorization"
	"forum-thread/pkg/api/forum-status"
	"forum-thread/pkg/api/forum-user"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
)

type IThreadService interface {
	CreateThread(ctx context.Context, threadName, threadDescription, threadColor string, allowedStatuses []string) (int, error)

	AllThreads(ctx context.Context) ([]*Thread, error)
}

type IThreadRepo interface {
	CreateThread(ctx context.Context, threadName, threadDescription, threadColor string, allowedStatuses []string) (int, error)

	AllThreads(ctx context.Context) ([]*Thread, error)
}

type ISubthreadService interface {
	CreateSubthread(ctx context.Context,
		threadID int,
		threadName,
		subthreadName,
		subthreadDescription string,
	) (int, error)

	AddViewToSubthread(ctx context.Context, subthreadID int) error
	AddViewToSubthreadPostprocessing(ctx context.Context, subthreadID int) error

	SubthreadsByThreadID(ctx context.Context, threadID int) ([]*Subthread, error)
}

type ISubthreadRepo interface {
	CreateSubthread(ctx context.Context,
		threadID int,
		threadName,
		subthreadName,
		subthreadDescription string,
	) (int, error)

	AddViewToSubthread(ctx context.Context, subthreadID int) error
	AddMessageCountToSubthread(ctx context.Context, subthreadID int) error
	UpdateSubthreadLastMessage(ctx context.Context, subthreadID int, subthreadLastMessageLogin, subthreadLastMessageText string) error

	SubthreadsByThreadID(ctx context.Context, threadID int) ([]*Subthread, error)
}

type ITopicService interface {
	CreateTopic(ctx context.Context,
		subthreadID,
		threadID,
		topicOwnerAccountID int,
		threadName,
		subthreadName,
		topicName,
		topicOwnerLogin string,
	) (int, error)
	CreateTopicPostprocessing(ctx context.Context, topicOwnerAccountID int) error

	AddViewToTopic(ctx context.Context, topicID int) error
	AddViewToTopicPostprocessing(ctx context.Context, topicID int) error

	CloseTopic(ctx context.Context,
		topicOwnerAccountID,
		adminAccountID,
		topicID int,
		topicName,
		adminLogin string,
	) error
	CloseTopicPostprocessing(ctx context.Context,
		topicOwnerAccountID,
		adminAccountID,
		topicID int,
		topicName,
		adminLogin string,
	) error
	RejectTopic(ctx context.Context, topicID int) error
	ApproveTopic(ctx context.Context, topicID int) error
	UpdateTopicAvatar(ctx context.Context, topicID int, fileName, extension string, fileBytes []byte) error
	DownloadTopicAvatar(ctx context.Context, fileUrl string) ([]byte, error)

	ChangeTopicPriority(ctx context.Context, subthreadID, topicID int, topicPriority int) error
	ChangeTopicPriorityPostprocessing(ctx context.Context, subthreadID, topicID int, topicPriority int) error

	TopicsBySubthreadID(ctx context.Context, subthreadID int) ([]*Topic, error)
	TopicsByAccountID(ctx context.Context, topicOwnerAccountID int) ([]*Topic, error)
	TopicsOnModeration(ctx context.Context) ([]*Topic, error)
}

type ITopicRepo interface {
	CreateTopic(ctx context.Context,
		subthreadID,
		threadID,
		topicOwnerAccountID int,
		subthreadName,
		threadName,
		topicName,
		topicOwnerLogin string,
		topicIsAuthor bool,
	) (int, error)

	UpdateTopicLastMessage(ctx context.Context, topicID int, topicLastMessageLogin, topicLastMessageText string) error
	AddMessageCountToTopic(ctx context.Context, topicID int) error
	AddViewToTopic(ctx context.Context, topicID int) error
	CloseTopic(ctx context.Context, topicID int) error
	ChangeTopicPriority(ctx context.Context, topicID, topicPriority int) error
	ApproveTopic(ctx context.Context, topicID int) error
	RejectTopic(ctx context.Context, topicID int) error
	UpdateTopicAvatar(ctx context.Context, topicID int, topicAvatarURL string) error
	UploadAvatar(ctx context.Context, file []byte, name string) (string, error)
	DownloadAvatar(ctx context.Context, fileURL string) ([]byte, error)

	TopicsBySubthreadID(ctx context.Context, subthreadID int) ([]*Topic, error)
	TopicsBySubthreadIDAndAccountID(ctx context.Context, subthreadID, accountID int) ([]*Topic, error)
	TopicsByAccountID(ctx context.Context, topicOwnerAccountID int) ([]*Topic, error)
	TopicsOnModeration(ctx context.Context) ([]*Topic, error)
	TopicsByAccountIDToday(ctx context.Context, accountID int) ([]*Topic, error)
}

type IMessageService interface {
	SendMessageToTopic(ctx context.Context,
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
	) ([]string, error)

	SendMessageToTopicPostprocessing(ctx context.Context,
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
	) (int, error)

	LikeMessage(ctx context.Context,
		topicID,
		messageOwnerAccountID,
		likerAccountID,
		likeMessageID,
		likeTypeID int,
		likerLogin,
		topicName,
		likeMessageText string,
	) error
	LikeMessagePostprocessing(ctx context.Context,
		topicID,
		messageOwnerAccountID,
		likerAccountID,
		likeMessageID,
		likeTypeID int,
		likerLogin,
		topicName,
		likeMessageText string,
	) error
	UnlikeMessage(ctx context.Context, likeMessageID, likerAccountID int) error
	ReportMessage(ctx context.Context, messageID, accountID int, reportText string) error
	ReportMessagePostprocessing(ctx context.Context, messageID, accountID int, reportText string) error
	EditMessage(ctx context.Context, messageID int, messageText string) error

	MessagesByText(ctx context.Context, text string) ([]*MessageSearch, error)
	MessagesByTopicID(ctx context.Context, accountID, topicID int) ([]*Message, []*Like, []*File, error)
	MessagesByAccountID(ctx context.Context, accountID int) ([]*Message, []*File, error)
	DownloadFile(ctx context.Context, fileURL string) ([]byte, error)
}

type IMessageRepo interface {
	CreateMessage(ctx context.Context,
		topicID,
		replyToMessageID,
		accountID int,
		login,
		text string,
	) (int, error)
	CreateMessageSearchIndex(ctx context.Context,
		topicID,
		messageID,
		accountID int,
		login,
		text string,
	) error
	CreateMessageLike(ctx context.Context, topicID, messageID, likeTypeID, accountID int) (int, error)
	DeleteMessageLike(ctx context.Context, likeMessageID, likerAccountID int) error

	IncrementLikeCountToMessage(ctx context.Context, messageID int) error
	DecrementLikeCountToMessage(ctx context.Context, messageID int) error
	AddReplyCountToMessage(ctx context.Context, messageID int) error
	AddReportCountToMessage(ctx context.Context, messageID int) error
	UploadFile(ctx context.Context, file []byte, name string) (string, error)
	AddFileToMessage(ctx context.Context, messageID, size int, url, name, extension string) error
	EditMessage(ctx context.Context, messageID int, messageText string) error

	MessagesByText(ctx context.Context, text string) ([]*MessageSearch, error)
	MessagesByTopicID(ctx context.Context, topicID int) ([]*Message, error)
	MessagesByAccountID(ctx context.Context, accountID int) ([]*Message, error)
	LikesByTopicIDAndAccountID(ctx context.Context, topicID, accountID int) ([]*Like, error)
	FilesByMessageID(ctx context.Context, messageID int) ([]*File, error)
	DownloadFile(ctx context.Context, fileURL string) ([]byte, error)

	CtxWithTx(ctx context.Context) (context.Context, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context)
}

type IAccountStatisticService interface {
	CreateAccountStatistic(ctx context.Context, accountID int) (int, error)
	StatisticByAccountID(ctx context.Context, accountID int) (*AccountStatistic, error)
}

type IAccountStatisticRepo interface {
	CreateAccountStatistic(ctx context.Context, accountID int) (int, error)
	AddSentMessagesToTopicsCount(ctx context.Context, accountID int) error
	AddCreatedTopicsCount(ctx context.Context, accountID int) error
	StatisticByAccountID(ctx context.Context, accountID int) ([]*AccountStatistic, error)
}

type IAuthorizationClient interface {
	CheckAuthorization(ctx context.Context, request echo.Context) (*forum_authorization.AuthorizationData, error)
}

type IStatusClient interface {
	StatusByAccountID(ctx context.Context, accountID int) ([]*forum_status.Status, error)
	AssignStatusToAccount(ctx context.Context, statusID, accountID int) error
	RevokeStatusFromAccount(ctx context.Context, statusID, accountID int) error
}

type ISupportClient interface {
	ReportMessage(ctx context.Context, accountID, messageID int, reportText string) error
}

type INotificationClient interface {
	NewMessageFromTopicNotification(ctx context.Context,
		topicOwnerAccountID,
		senderMessageID,
		senderAccountID,
		topicID int,
		senderMessageText,
		TopicName,
		senderLogin string,
	) error
	NewMentionFromTopicNotification(ctx context.Context,
		mentionedAccountID,
		senderMessageID,
		senderAccountID int,
		senderMessageText,
		topicName,
		senderLogin string,
	) error
	NewLikeMessageFromTopicNotification(ctx context.Context,
		messageOwnerAccountID,
		senderMessageID,
		likerAccountID,
		topicID int,
		senderMessageText,
		topicName,
		likerLogin string,
	) error
	NewMessageReplyFromTopicNotification(ctx context.Context,
		replyMessageOwnerAccountID,
		senderMessageID,
		senderAccountID,
		topicID int,
		senderMessageText,
		topicName,
		senderLogin string,
	) error
	TopicClosedNotificationRequest(ctx context.Context,
		topicOwnerAccountID,
		adminAccountID,
		topicID int,
		topicName,
		adminLogin string,
	) error
}

type IUserClient interface {
	UserByLogin(ctx context.Context, login string) (*forum_user.GetUserByLoginResponse, error)
}

type IDatabase interface {
	Insert(ctx context.Context, query string, args ...any) (int, error)
	Select(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Delete(ctx context.Context, query string, args ...any) error
	Update(ctx context.Context, query string, args ...any) error

	CtxWithTx(ctx context.Context) (context.Context, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context)

	CreateTable(ctx context.Context, query string) error
	DropTable(ctx context.Context, query string) error
}

type IFullTextSearchEngine interface {
	CreateIndexes(indexes []string) error
	AddDocuments(indexName string, documents []any) error
	SimpleSearch(indexName, query string) ([]byte, error)
}

type IMessageBroker interface {
	Publish(ctx context.Context, queueName string, body any) error
	Subscribe(queueName string, handler func(event amqp091.Delivery) error)
}

type IWeedFS interface {
	Upload(file []byte, filename string, size int64, collection string) (string, error)
	Download(FileID string) ([]byte, error)
	Update(file []byte, fileID string, filename string, size int64, collection string) error
}
