package unit

import (
	"context"
	"forum-authentication/infrastructure/pg"
	"forum-authentication/internal/app/server"
	"forum-authentication/internal/config"
	"forum-authentication/internal/model"
)

import (
	"forum-authentication/test/mocks/forum-admin"
	"forum-authentication/test/mocks/forum-authorization"
	"forum-authentication/test/mocks/forum-thread"
	"forum-authentication/test/mocks/forum-user"
)

import (
	AuthenticationRepo "forum-authentication/internal/repo/authentication"
	AuthenticationService "forum-authentication/internal/service/authentication"
)

func New() *TestConfig {
	cfg := config.New()
	db := pg.New(
		"postgres",
		"postgres",
		"forum-authentication-db-unit-test",
		"5432",
		"postgres",
	)

	authorizationClient := forum_authorization.New()
	threadClient := forum_thread.New()
	userClient := forum_user.New()
	adminClient := forum_admin.New()

	authenticationRepo := AuthenticationRepo.New(db)
	authenticationService := AuthenticationService.New(
		authenticationRepo,
		authorizationClient,
		userClient,
		adminClient,
		threadClient,
		cfg.Forum.PasswordSecretKey,
	)

	authenticationClient := NewAuthenticationClient("localhost", "8003", cfg.Forum.InterServerSecretKey)

	go server.Run(
		db,
		authenticationService,
		authorizationClient,
		"8003",
		"localhost",
		cfg.Forum.InterServerSecretKey,
	)

	return &TestConfig{
		db:                    db,
		authenticationService: authenticationService,
		authenticationRepo:    authenticationRepo,
		authenticationClient:  authenticationClient,
	}
}

type TestConfig struct {
	db                    model.IDatabase
	authenticationService model.IAuthenticationService
	authenticationRepo    model.IAuthenticationRepo
	authenticationClient  *ClientAuthentication
	interServerSecretKey  string
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
