package unit

import (
	"context"
	"forum-dialog/infrastructure/pg"
	"forum-dialog/infrastructure/redis"
	"forum-dialog/infrastructure/weed"
	"forum-dialog/internal/app/server"
	"forum-dialog/internal/config"
	"forum-dialog/internal/model"
)

import (
	"forum-dialog/test/mocks/forum-authorization"
	"forum-dialog/test/mocks/rabbitmq"
)

import (
	DialogRepo "forum-dialog/internal/repo/dialog"
	DialogService "forum-dialog/internal/service/dialog"
)

func New() *TestConfig {
	cfg := config.New()
	db := pg.New(
		"postgres",
		"postgres",
		"forum-dialog-db-unit-test",
		"5432",
		"postgres",
	)
	redisClient := redis.New(
		"forum-dialog-redis-unit-test",
		"6379",
		"password111",
		1,
	)
	messageBroker := rabbitmq.New()
	weedFS := weed.New("http://forum-dialog-weed-filer-unit-test:8888", "http://forum-dialog-weed-master-unit-test:9333")

	wsConnManager := server.NewWsConnManager(cfg.PodID, redisClient, messageBroker)
	authorizationClient := forum_authorization.New()

	dialogRepo := DialogRepo.New(db, weedFS)
	dialogService := DialogService.New(dialogRepo, wsConnManager)

	dialogClient := NewDialogClient("localhost", "8003")

	go server.Run(
		db,
		dialogService,
		authorizationClient,
		wsConnManager,
		messageBroker,
		"8003",
		cfg.PodID,
	)

	return &TestConfig{
		db:            db,
		dialogService: dialogService,
		dialogClient:  dialogClient,
	}
}

type TestConfig struct {
	db            model.IDatabase
	dialogService model.IDialogService
	dialogClient  *ClientDialog
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
