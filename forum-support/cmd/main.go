package cmd

import (
	"forum-support/infrastructure/pg"
	"forum-support/infrastructure/rabbitmq"
	"forum-support/infrastructure/redis"
	"forum-support/internal/app/server"
	"forum-support/internal/config"
	DialogRepo "forum-support/internal/repo/dialog"
	SupportRequestRepo "forum-support/internal/repo/support-request"
	DialogService "forum-support/internal/service/dialog"
	SupportRequestService "forum-support/internal/service/support-request"
	"forum-support/pkg/api/forum-authorization"
	"forum-support/pkg/api/forum-notifications"
	"forum-support/pkg/logger"
)

func main() {
	logger.New()
	cfg := config.New()

	db := pg.New(cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Dbname)
	redisClient := redis.New(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, 1)
	messageBroker := rabbitmq.New(
		cfg.Rabbitmq.Host,
		cfg.Rabbitmq.Port,
		cfg.Rabbitmq.Username,
		cfg.Rabbitmq.Password,
	)
	wsConnManager := server.NewWsConnManager(cfg.PodID, redisClient, messageBroker)

	supportRequestRepo := SupportRequestRepo.New(db)
	notificationClient := forum_notifications.New(cfg.Forum.Notification.Host, cfg.Forum.Notification.Port)
	supportRequestService := SupportRequestService.New(supportRequestRepo, notificationClient)

	dialogRepo := DialogRepo.New(db)
	dialogService := DialogService.New(dialogRepo, wsConnManager)

	authorizationClient := forum_authorization.New(cfg.Forum.Authorization.Host, cfg.Forum.Authorization.Port)

	server.Run(
		db,
		supportRequestService,
		dialogService,
		authorizationClient,
		wsConnManager,
		cfg.Forum.Notification.Port,
	)
}
