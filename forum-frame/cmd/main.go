package main

import (
	"forum-frame/infrastructure/pg"
	"forum-frame/infrastructure/weed"
	"forum-frame/pkg/api/forum-authorization"
	"forum-frame/pkg/api/forum-payment"
	"forum-frame/pkg/logger"
)

import (
	"forum-frame/internal/app/server"
	"forum-frame/internal/config"
	FrameRepo "forum-frame/internal/repo/frame"
	FrameService "forum-frame/internal/service/frame"
)

func main() {
	logger.New()
	cfg := config.New()

	db := pg.New(cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Dbname)
	weedFS := weed.New(cfg.WeedFS.FilerUrl, cfg.WeedFS.MasterUrl)

	authorizationClient := forum_authorization.New(cfg.Forum.Authorization.Host, cfg.Forum.Authorization.Port)
	paymentClient := forum_payment.New(cfg.Forum.Payment.Host, cfg.Forum.Payment.Port)

	frameRepo := FrameRepo.New(db, weedFS)
	frameService := FrameService.New(frameRepo, paymentClient)

	server.Run(
		db,
		frameService,
		authorizationClient,
		cfg.Forum.InterServerSecretKey,
		cfg.Forum.Frame.Port,
	)
}
