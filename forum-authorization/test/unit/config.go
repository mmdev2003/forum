package unit

import (
	"context"
	"forum-authorization/infrastructure/pg"
	"forum-authorization/internal/app/server"
	"forum-authorization/internal/config"
	"forum-authorization/internal/model"
)

import (
	AuthorizationRepo "forum-authorization/internal/repo/authorization"
	AuthorizationService "forum-authorization/internal/service/authorization"
)

func New() *TestConfig {
	cfg := config.New()
	db := pg.New(
		"postgres",
		"postgres",
		"forum-authorization-db-unit-test",
		"5432",
		"postgres",
	)

	authorizationRepo := AuthorizationRepo.New(db)
	authorizationService := AuthorizationService.New(authorizationRepo, cfg.Forum.JwtSecretKey)

	authorizationClient := NewAuthorizationClient("localhost", "8003")

	go server.Run(
		db,
		authorizationService,
		"8003",
		"localhost",
	)

	return &TestConfig{
		db:                   db,
		authorizationService: authorizationService,
		authorizationRepo:    authorizationRepo,
		authorizationClient:  authorizationClient,
	}
}

type TestConfig struct {
	db                   model.IDatabase
	authorizationService model.IAuthorizationService
	authorizationRepo    model.IAuthorizationRepo
	authorizationClient  *ClientAuthorization
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
