package main

import (
	"forum-authorization/infrastructure/pg"
	"forum-authorization/pkg/logger"
)

import (
	"forum-authorization/internal/app/server"
	"forum-authorization/internal/config"
	AuthorizationRepo "forum-authorization/internal/repo/authorization"
	AuthorizationService "forum-authorization/internal/service/authorization"
)

func main() {
	logger.New()
	cfg := config.New()

	db := pg.New(cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Dbname)

	authorizationRepo := AuthorizationRepo.New(db)
	authorizationService := AuthorizationService.New(authorizationRepo, cfg.Forum.JwtSecretKey)

	server.Run(
		db,
		authorizationService,
		cfg.Forum.Authorization.Port,
		cfg.Forum.Domain,
	)
}
