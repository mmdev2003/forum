package unit

import (
	"context"
	"forum-admin/infrastructure/pg"
	"forum-admin/internal/app/server"
	"forum-admin/internal/config"
	"forum-admin/internal/model"
)

import (
	AdminRepo "forum-admin/internal/repo/admin"
	AdminService "forum-admin/internal/service/admin"
)

func New() *TestConfig {
	cfg := config.New()
	db := pg.New(
		"postgres",
		"postgres",
		"forum-admin-db-unit-test",
		"5432",
		"postgres",
	)
	adminRepo := AdminRepo.New(db)
	adminService := AdminService.New(adminRepo)

	adminClient := NewAdminClient("localhost", "8003", cfg.Forum.InterServerSecretKey)

	go server.Run(
		db,
		adminService,
		cfg.Forum.InterServerSecretKey,
		"8003",
	)

	return &TestConfig{
		db:                   db,
		adminService:         adminService,
		adminClient:          adminClient,
		interServerSecretKey: cfg.Forum.InterServerSecretKey,
	}
}

type TestConfig struct {
	db                   model.IDatabase
	adminService         model.IAdminService
	adminClient          *ClientAdmin
	interServerSecretKey string
}

func (testConfig *TestConfig) PrepareDB() {
	err := testConfig.db.DropTable(context.Background(), model.DropTableQuery)
	if err != nil {
		panic(err)
	}

	err = testConfig.db.CreateTable(context.Background(), model.CreateTableQuery)
	if err != nil {
		panic(err)
	}
}

var testConfig *TestConfig = New()
