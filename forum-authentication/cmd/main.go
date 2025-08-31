package main

import (
	"forum-authentication/infrastructure/pg"
	"forum-authentication/internal/app/server"
	"forum-authentication/internal/config"
	"forum-authentication/pkg/logger"
)

import (
	"forum-authentication/pkg/api/forum-admin"
	"forum-authentication/pkg/api/forum-authorization"
	"forum-authentication/pkg/api/forum-thread"
	"forum-authentication/pkg/api/forum-user"
)

import (
	AuthenticationRepo "forum-authentication/internal/repo/authentication"
	AuthenticationService "forum-authentication/internal/service/authentication"
)

func main() {
	logger.New()
	cfg := config.New()

	db := pg.New(cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Dbname)

	authorizationClient := forum_authorization.New(cfg.Forum.Authorization.Host, cfg.Forum.Authorization.Port)
	userClient := forum_user.New(cfg.Forum.User.Host, cfg.Forum.User.Port)
	adminClient := forum_admin.New(cfg.Forum.Admin.Host, cfg.Forum.Admin.Port, cfg.Forum.InterServerSecretKey)
	threadClient := forum_thread.New(cfg.Forum.Thread.Host, cfg.Forum.Thread.Port)

	authenticationRepo := AuthenticationRepo.New(db)
	authenticationService := AuthenticationService.New(
		authenticationRepo,
		authorizationClient,
		userClient,
		adminClient,
		threadClient,
		cfg.Forum.PasswordSecretKey,
	)

	server.Run(
		db,
		authenticationService,
		authorizationClient,
		cfg.Forum.Authentication.Port,
		cfg.Forum.Domain,
		cfg.Forum.InterServerSecretKey,
	)
}
