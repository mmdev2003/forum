package unit

import (
	"context"
	"forum-status/infrastructure/pg"
	"forum-status/internal/app/server"
	"forum-status/internal/config"
	"forum-status/internal/model"
)

import (
	"forum-status/test/mocks/forum-authorization"
	"forum-status/test/mocks/forum-payment"
)

import (
	StatusRepo "forum-status/internal/repo/status"
	StatusService "forum-status/internal/service/status"
)

func New() *TestConfig {
	cfg := config.New()
	db := pg.New(
		"postgres",
		"postgres",
		"forum-status-db-unit-test",
		"5432",
		"postgres",
	)

	authorizationClient := forum_authorization.New()
	paymentClient := forum_payment.New()

	statusRepo := StatusRepo.New(db)
	statusService := StatusService.New(statusRepo, paymentClient)

	statusClient := NewStatusClient("localhost", "8003")

	go server.Run(
		db,
		statusService,
		authorizationClient,
		cfg.Forum.InterServerSecretKey,
		"8003",
	)

	return &TestConfig{
		db:                   db,
		statusService:        statusService,
		statusClient:         statusClient,
		interServerSecretKey: cfg.Forum.InterServerSecretKey,
	}
}

type TestConfig struct {
	db                   model.IDatabase
	statusService        model.IStatusService
	statusClient         *ClientStatus
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
