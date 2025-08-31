package main

import (
	"forum-dialog/infrastructure/pg"
	"forum-dialog/infrastructure/rabbitmq"
	"forum-dialog/infrastructure/redis"
	"forum-dialog/infrastructure/weed"
	"forum-dialog/pkg/api/forum-authorization"
	"forum-dialog/pkg/logger"
)

import (
	"forum-dialog/internal/app/server"
	"forum-dialog/internal/config"
	DialogRepo "forum-dialog/internal/repo/dialog"
	DialogService "forum-dialog/internal/service/dialog"
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
	redisClient := redis.New(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		1,
	)
	messageBroker := rabbitmq.New(
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
	)
	weedFS := weed.New(cfg.WeedFS.FilerUrl, cfg.WeedFS.MasterUrl)

	wsConnManager := server.NewWsConnManager(cfg.PodID, redisClient, messageBroker)

	dialogRepo := DialogRepo.New(db, weedFS)
	dialogService := DialogService.New(dialogRepo, wsConnManager)

	authorizationClient := forum_authorization.New(cfg.Forum.Authorization.Host, cfg.Forum.Authorization.Port)

	server.Run(
		db,
		dialogService,
		authorizationClient,
		wsConnManager,
		messageBroker,
		cfg.Forum.Dialog.Port,
		cfg.PodID,
	)
}
