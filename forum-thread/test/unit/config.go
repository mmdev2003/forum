package unit

import (
	"context"
	"forum-thread/infrastructure/meilisearch"
	"forum-thread/infrastructure/pg"
	"forum-thread/infrastructure/weed"
	"forum-thread/internal/app/server"
	"forum-thread/internal/model"
	"forum-thread/test/mocks/forum-user"
)

import (
	"forum-thread/test/mocks/forum-authorization"
	"forum-thread/test/mocks/forum-notification"
	"forum-thread/test/mocks/forum-status"
	"forum-thread/test/mocks/forum-support"
	"forum-thread/test/mocks/rabbitmq"
)

import (
	AccountStatisticRepo "forum-thread/internal/repo/account-statistic"
	MessageRepo "forum-thread/internal/repo/message"
	SubthreadRepo "forum-thread/internal/repo/subthread"
	ThreadRepo "forum-thread/internal/repo/thread"
	TopicRepo "forum-thread/internal/repo/topic"
	AccountStatisticService "forum-thread/internal/service/account-statistic"
	MessageService "forum-thread/internal/service/message"
	SubthreadService "forum-thread/internal/service/subthread"
	ThreadService "forum-thread/internal/service/thread"
	TopicService "forum-thread/internal/service/topic"
)

func New() *TestConfig {
	db := pg.New(
		"postgres",
		"postgres",
		"forum-thread-db-unit-test",
		"5432",
		"postgres",
	)

	messageBroker := rabbitmq.New()

	weedFS := weed.New("http://forum-thread-weed-filer-unit-test:8888", "http://forum-thread-weed-master-unit-test:9333")

	fullTextSearchEngine := meilisearch.New(
		"forum-meilisearch-thread-unit-test",
		"7700",
		"qwertyuiop",
	)
	indexes := []string{
		model.MessageFullTextSearchIndex,
	}
	err := fullTextSearchEngine.CreateIndexes(indexes)
	if err != nil {
		panic(err)
	}

	supportClient := forum_support.New()
	authorizationClient := forum_authorization.New()
	notificationClient := forum_notification.New()
	userClient := forum_user.New()
	statusClient := forum_status.New()

	threadRepo := ThreadRepo.New(db, fullTextSearchEngine)
	subthreadRepo := SubthreadRepo.New(db, fullTextSearchEngine)
	topicRepo := TopicRepo.New(db, fullTextSearchEngine)
	messageRepo := MessageRepo.New(db, fullTextSearchEngine, weedFS)
	accountStatisticRepo := AccountStatisticRepo.New(db)

	threadService := ThreadService.New(threadRepo, messageBroker)
	subthreadService := SubthreadService.New(subthreadRepo, messageBroker)
	topicService := TopicService.New(topicRepo, accountStatisticRepo, messageBroker, notificationClient, statusClient)
	messageService := MessageService.New(
		subthreadRepo,
		topicRepo,
		messageRepo,
		accountStatisticRepo,
		supportClient,
		statusClient,
		notificationClient,
		userClient,
		messageBroker,
	)
	accountStatisticService := AccountStatisticService.New(accountStatisticRepo)

	threadClient := NewThreadClient("localhost", "8003")

	go server.Run(
		db,
		fullTextSearchEngine,
		messageBroker,
		threadService,
		subthreadService,
		topicService,
		messageService,
		accountStatisticService,
		authorizationClient,
		"8003",
	)

	return &TestConfig{
		db:                      db,
		threadService:           threadService,
		threadRepo:              threadRepo,
		subthreadService:        subthreadService,
		subthreadRepo:           subthreadRepo,
		topicService:            topicService,
		topicRepo:               topicRepo,
		messageService:          messageService,
		messageRepo:             messageRepo,
		accountStatisticService: accountStatisticService,
		accountStatisticRepo:    accountStatisticRepo,
		threadClient:            threadClient,
	}
}

type TestConfig struct {
	db                      model.IDatabase
	threadService           model.IThreadService
	threadRepo              model.IThreadRepo
	subthreadService        model.ISubthreadService
	subthreadRepo           model.ISubthreadRepo
	topicService            model.ITopicService
	topicRepo               model.ITopicRepo
	messageService          model.IMessageService
	messageRepo             model.IMessageRepo
	accountStatisticService model.IAccountStatisticService
	accountStatisticRepo    model.IAccountStatisticRepo
	threadClient            *ClientThread
}

func (testConfig *TestConfig) PrepareDB() {
	err := testConfig.db.DropTable(context.Background(), model.DropTableQuery)
	if err != nil {
		panic(err)
	}

	err = testConfig.db.CreateTable(context.Background(), model.CreateTableQuery)
	if err != nil {
		panic(err)
	}
}

var testConfig *TestConfig = New()
