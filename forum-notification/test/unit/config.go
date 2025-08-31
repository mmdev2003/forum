package unit

import (
	"context"
	"forum-notification/infrastructure/pg"
	"forum-notification/internal/app/server"
	"forum-notification/internal/config"
	"forum-notification/internal/model"
	NotificationFilterRepo "forum-notification/internal/repo/filter"
	NotificationRepo "forum-notification/internal/repo/notification"
	NotificationFilterService "forum-notification/internal/service/filter"
	NotificationService "forum-notification/internal/service/notification"
	forumauthorization "forum-notification/test/mocks/forum-authorization"
)

func New() *TestConfig {
	db := pg.New(
		"postgres",
		"postgres",
		"forum-notification-db-unit-test",
		"5432",
		"postgres",
	)

	cfg := config.New()

	authorizationClient := forumauthorization.New()
	wsConnManager := server.NewWsConnManager()

	notificationRepo := NotificationRepo.New(db)
	notificationFilterRepo := NotificationFilterRepo.New(db)
	notificationFilterService := NotificationFilterService.New(notificationFilterRepo)
	notificationService := NotificationService.New(notificationRepo, notificationFilterService, wsConnManager)

	notificationClient := NewNotificationClient("localhost", "8002")

	go server.Run(
		db,
		notificationService,
		notificationFilterService,
		wsConnManager,
		authorizationClient,
		cfg.Forum.Notification.Port,
	)

	return &TestConfig{
		db:                  db,
		notificationService: notificationService,
		notificationClient:  notificationClient,
	}
}

type TestConfig struct {
	db                  model.IDatabase
	notificationService model.INotification
	notificationClient  *NotificationClient
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
