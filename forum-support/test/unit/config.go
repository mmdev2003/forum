package unit

import (
	"context"
	"forum-support/infrastructure/pg"
	"forum-support/infrastructure/redis"
	"forum-support/internal/app/server"
	"forum-support/internal/config"
	"forum-support/internal/model"
	DialogRepo "forum-support/internal/repo/dialog"
	DialogService "forum-support/internal/service/dialog"

	forum_authorization "forum-support/test/mocks/forum-authorization"

	forum_notification "forum-support/test/mocks/forum-notification"
	"forum-support/test/mocks/rabbitmq"

	SupportRequestRepo "forum-support/internal/repo/support-request"

	SupportRequestService "forum-support/internal/service/support-request"
)

func New() *TestConfig {
	cfg := config.New()

	db := pg.New(
		"postgres",
		"postgres",
		"forum-support-db-unit-test",
		"5432",
		"postgres",
	)
	redisClient := redis.New(
		"forum-support-redis-unit-test",
		"6379",
		"password111",
		1,
	)
	messageBroker := rabbitmq.New()

	wsConnManager := server.NewWsConnManager(cfg.PodID, redisClient, messageBroker)
	authorizationClient := forum_authorization.New()
	notificationClient := forum_notification.New()

	supportRequestRepo := SupportRequestRepo.New(db)
	supportRequestService := SupportRequestService.New(supportRequestRepo, notificationClient)

	dialogRepo := DialogRepo.New(db)
	dialogService := DialogService.New(dialogRepo, wsConnManager)

	supportClient := NewSupportClient("localhost", "8003")
	dialogClient := NewDialogClient("localhost", "8003")

	go server.Run(
		db,
		supportRequestService,
		dialogService,
		authorizationClient,
		wsConnManager,
		"8003",
	)

	return &TestConfig{
		db:            db,
		SupportClient: supportClient,
		DialogClient:  dialogClient,
	}
}

type TestConfig struct {
	db            model.IDatabase
	SupportClient *SupportClient
	DialogClient  *ClientDialog
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
