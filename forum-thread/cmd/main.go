package main

import (
	"forum-thread/infrastructure/meilisearch"
	"forum-thread/infrastructure/pg"
	"forum-thread/infrastructure/rabbitmq"
	"forum-thread/infrastructure/weed"
	"forum-thread/internal/config"
	"forum-thread/pkg/logger"
)

import (
	"forum-thread/pkg/api/forum-authorization"
	"forum-thread/pkg/api/forum-notification"
	"forum-thread/pkg/api/forum-status"
	"forum-thread/pkg/api/forum-support"
	"forum-thread/pkg/api/forum-user"
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

import (
	"forum-thread/internal/app/server"
)

func main() {
	logger.New()
	cfg := config.New()

	db := pg.New(
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Dbname,
	)

	fullTextSearchEngine := meilisearch.New(
		cfg.Meilisearch.Host,
		cfg.Meilisearch.Port,
		cfg.Meilisearch.ApiKey,
	)
	messageBroker := rabbitmq.New(
		cfg.Rabbitmq.Host,
		cfg.Rabbitmq.Port,
		cfg.Rabbitmq.Username,
		cfg.Rabbitmq.Password,
	)
	weedFS := weed.New(cfg.WeedFS.FilerUrl, cfg.WeedFS.MasterUrl)

	supportClient := forum_support.New(cfg.Forum.Support.Host, cfg.Forum.Support.Port)
	authorizationClient := forum_authorization.New(cfg.Forum.Authorization.Host, cfg.Forum.Authorization.Port)
	notificationClient := forum_notification.New(cfg.Forum.Notification.Host, cfg.Forum.Notification.Port)
	userClient := forum_user.New(cfg.Forum.User.Host, cfg.Forum.User.Port)
	statusClient := forum_status.New(cfg.Forum.Status.Host, cfg.Forum.Status.Port, cfg.Forum.InterServerSecretKey)

	threadRepo := ThreadRepo.New(db, fullTextSearchEngine)
	subthreadRepo := SubthreadRepo.New(db, fullTextSearchEngine)
	topicRepo := TopicRepo.New(db, fullTextSearchEngine, weedFS)
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

	server.Run(
		db,
		fullTextSearchEngine,
		messageBroker,
		threadService,
		subthreadService,
		topicService,
		messageService,
		accountStatisticService,
		authorizationClient,
		cfg.Forum.Thread.Port,
	)

}
