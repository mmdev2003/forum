package main

import (
	"forum-admin/infrastructure/pg"
	"forum-admin/pkg/logger"
)

import (
	"forum-admin/internal/app/server"
	"forum-admin/internal/config"
	FrameRepo "forum-admin/internal/repo/admin"
	FrameService "forum-admin/internal/service/admin"
)

func main() {
	logger.New()
	cfg := config.New()

	db := pg.New(cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Dbname)

	adminRepo := FrameRepo.New(db)
	adminService := FrameService.New(adminRepo)

	server.Run(
		db,
		adminService,
		cfg.Forum.InterServerSecretKey,
		cfg.Forum.Admin.Port,
	)
}
