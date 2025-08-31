package main

import (
	"flag"
	"forum-payment/infrastructure/pg"
	"forum-payment/internal/app/btc_confirmator"
	"forum-payment/internal/app/server"
	"forum-payment/internal/config"
	"forum-payment/pkg/api/forum-authorization"
	"forum-payment/pkg/logger"
)

import (
	"forum-payment/pkg/api/forum-frame"
	"forum-payment/pkg/api/forum-status"
)

import (
	PaymentRepo "forum-payment/internal/repo/payment"
	PaymentService "forum-payment/internal/service/payment"
)

func main() {
	app := flag.String("app", "", "Выбор приложения")
	flag.Parse()

	logger.New()
	cfg := config.New()

	db := pg.New(cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Dbname)

	authorizationClient := forum_authorization.New(cfg.Forum.Authorization.Host, cfg.Forum.Authorization.Port)
	frameClient := forum_frame.New(cfg.Forum.Frame.Host, cfg.Forum.Frame.Port, cfg.Forum.InterServerSecretKey)
	statusClient := forum_status.New(cfg.Forum.Status.Host, cfg.Forum.Status.Port, cfg.Forum.InterServerSecretKey)

	paymentRepo := PaymentRepo.New(db)
	paymentService := PaymentService.New(paymentRepo, frameClient, statusClient, cfg.Forum.BtcAddress)

	if *app == "server" {
		server.Run(
			db,
			paymentService,
			authorizationClient,
			cfg.Forum.Payment.Port,
		)
	}
	if *app == "confirmator" {
		btc_confirmator.Run(paymentService, paymentRepo, cfg.Forum.BtcAddress)
	}
}
