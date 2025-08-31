package unit

import (
	"context"
	"forum-user/infrastructure/meilisearch"
	"forum-user/infrastructure/pg"
	"forum-user/infrastructure/weed"
	"forum-user/internal/app/server"
	"forum-user/internal/model"
)

import (
	"forum-user/test/mocks/forum-authorization"
	"forum-user/test/mocks/forum-notification"
)

import (
	UserRepo "forum-user/internal/repo/user"
	UserService "forum-user/internal/service/user"
)

func New() *TestConfig {
	db := pg.New(
		"postgres",
		"postgres",
		"forum-user-db-unit-test",
		"5432",
		"postgres",
	)
	weedFS := weed.New("http://forum-user-weed-filer-unit-test:8888", "http://forum-user-weed-master-unit-test:9333")
	fullTextSearchEngine := meilisearch.New(
		"forum-meilisearch-user-unit-test",
		"7700",
		"qwertyuiop",
	)
	indexes := []string{
		model.UserFullTextSearchIndex,
	}
	err := fullTextSearchEngine.CreateIndexes(indexes)
	if err != nil {
		panic(err)
	}

	authorizationClient := forum_authorization.New()
	notificationClient := forum_notification.New()

	userRepo := UserRepo.New(db, fullTextSearchEngine, weedFS)
	userService := UserService.New(userRepo, notificationClient)

	userClient := NewUserClient("localhost", "8003")

	go server.Run(
		db,
		fullTextSearchEngine,
		userService,
		authorizationClient,
		"8003",
	)

	return &TestConfig{
		db:          db,
		userService: userService,
		userClient:  userClient,
	}
}

type TestConfig struct {
	db          model.IDatabase
	userService model.IUserService
	userClient  *ClientUser
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
