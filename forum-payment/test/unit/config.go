package unit

import (
	"context"
	"forum-payment/infrastructure/pg"
	"forum-payment/internal/app/server"
	"forum-payment/internal/config"
	"forum-payment/internal/model"
)

import (
	"forum-payment/test/mocks/forum-authorization"
	"forum-payment/test/mocks/forum-frame"
	"forum-payment/test/mocks/forum-status"
)

import (
	PaymentRepo "forum-payment/internal/repo/payment"
	PaymentService "forum-payment/internal/service/payment"
)

func New() *TestConfig {
	cfg := config.New()
	db := pg.New(
		"postgres",
		"postgres",
		"forum-payment-db-unit-test",
		"5432",
		"postgres",
	)

	authorizationClient := forum_authorization.New()
	frameClient := forum_frame.New()
	statusClient := forum_status.New()

	paymentRepo := PaymentRepo.New(db)
	paymentService := PaymentService.New(paymentRepo, frameClient, statusClient, cfg.Forum.InterServerSecretKey)

	paymentClient := NewPaymentClient("localhost", "8003")

	go server.Run(
		db,
		paymentService,
		authorizationClient,
		"8003",
	)

	return &TestConfig{
		db:                   db,
		paymentService:       paymentService,
		paymentRepo:          paymentRepo,
		paymentClient:        paymentClient,
		interServerSecretKey: cfg.Forum.InterServerSecretKey,
	}
}

type TestConfig struct {
	db                   model.IDatabase
	paymentService       model.IPaymentService
	paymentRepo          model.IPaymentRepo
	paymentClient        *ClientPayment
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
