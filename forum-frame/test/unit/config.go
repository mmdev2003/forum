package unit

import (
	"context"
	"forum-frame/infrastructure/pg"
	"forum-frame/infrastructure/weed"
	"forum-frame/internal/app/server"
	"forum-frame/internal/config"
	"forum-frame/internal/model"
)

import (
	"forum-frame/test/mocks/forum-authorization"
	"forum-frame/test/mocks/forum-payment"
)

import (
	FrameRepo "forum-frame/internal/repo/frame"
	FrameService "forum-frame/internal/service/frame"
)

func New() *TestConfig {
	cfg := config.New()
	db := pg.New(
		"postgres",
		"postgres",
		"forum-frame-db-unit-test",
		"5432",
		"postgres",
	)
	weedFS := weed.New("http://forum-frame-weed-filer-unit-test:8888", "http://forum-frame-weed-master-unit-test:9333")

	authorizationClient := forum_authorization.New()
	paymentClient := forum_payment.New()

	frameRepo := FrameRepo.New(db, weedFS)
	frameService := FrameService.New(frameRepo, paymentClient)

	frameClient := NewFrameClient("localhost", "8003")

	go server.Run(
		db,
		frameService,
		authorizationClient,
		cfg.Forum.InterServerSecretKey,
		"8003",
	)

	return &TestConfig{
		db:                   db,
		frameService:         frameService,
		frameClient:          frameClient,
		interServerSecretKey: cfg.Forum.InterServerSecretKey,
	}
}

type TestConfig struct {
	db                   model.IDatabase
	frameService         model.IFrameService
	frameClient          *ClientFrame
	interServerSecretKey string
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
