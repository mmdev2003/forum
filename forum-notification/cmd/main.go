package main

import (
	"forum-notification/infrastructure/pg"
	"forum-notification/internal/app/server"
	"forum-notification/internal/config"
	NotificationFilterRepo "forum-notification/internal/repo/filter"
	NotificationRepo "forum-notification/internal/repo/notification"
	NotificationFilterService "forum-notification/internal/service/filter"
	NotificationService "forum-notification/internal/service/notification"
	forum_authorization "forum-notification/pkg/api/forum-authorization"
	"forum-notification/pkg/logger"
)

func main() {
	logger.New()
	cfg := config.New()

	db := pg.New(cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Dbname)
	wsConnManager := server.NewWsConnManager()

	notificationRepo := NotificationRepo.New(db)
	notificationFilterRepo := NotificationFilterRepo.New(db)
	notificationFilterService := NotificationFilterService.New(notificationFilterRepo)
	notificationService := NotificationService.New(notificationRepo, notificationFilterService, wsConnManager)

	authorizationClient := forum_authorization.New(cfg.Forum.Authorization.Host, cfg.Forum.Authorization.Port)

	server.Run(
		db,
		notificationService,
		notificationFilterService,
		wsConnManager,
		authorizationClient,
		cfg.Forum.Notification.Port,
	)
}
