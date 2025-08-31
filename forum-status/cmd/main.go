package main

import (
	"forum-status/infrastructure/pg"
	"forum-status/pkg/api/forum-authorization"
	"forum-status/pkg/api/forum-payment"
	"forum-status/pkg/logger"
)

import (
	"forum-status/internal/app/server"
	"forum-status/internal/config"
	StatusRepo "forum-status/internal/repo/status"
	StatusService "forum-status/internal/service/status"
)

func main() {
	logger.New()
	cfg := config.New()

	db := pg.New(cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Dbname)

	authorizationClient := forum_authorization.New(cfg.Forum.Authorization.Host, cfg.Forum.Authorization.Port)
	paymentClient := forum_payment.New(cfg.Forum.Payment.Host, cfg.Forum.Payment.Port)

	statusRepo := StatusRepo.New(db)
	statusService := StatusService.New(statusRepo, paymentClient)

	server.Run(
		db,
		statusService,
		authorizationClient,
		cfg.Forum.InterServerSecretKey,
		cfg.Forum.Status.Port,
	)
}
