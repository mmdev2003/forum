package main

import (
	"forum-user/infrastructure/meilisearch"
	"forum-user/infrastructure/pg"
	"forum-user/infrastructure/weed"
	"forum-user/internal/config"
	"forum-user/pkg/api/forum-authorization"
	"forum-user/pkg/api/forum-notification"
	"forum-user/pkg/logger"
)

import (
	UserRepo "forum-user/internal/repo/user"
	UserService "forum-user/internal/service/user"
)

import (
	"forum-user/internal/app/server"
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
	weedFS := weed.New(cfg.WeedFS.FilerUrl, cfg.WeedFS.MasterUrl)

	authorizationClient := forum_authorization.New(cfg.Forum.Authorization.Host, cfg.Forum.Authorization.Port)
	notificationClient := forum_notification.New(cfg.Forum.Notification.Host, cfg.Forum.Notification.Port)

	userRepo := UserRepo.New(db, fullTextSearchEngine, weedFS)
	userService := UserService.New(userRepo, notificationClient)

	server.Run(
		db,
		fullTextSearchEngine,
		userService,
		authorizationClient,
		cfg.Forum.User.Port,
	)

}
