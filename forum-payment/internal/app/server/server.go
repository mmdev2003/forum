package server

import (
	"context"
	"forum-payment/internal/model"
	"forum-payment/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

import (
	paymentHttpHandler "forum-payment/internal/controller/http/handler/payment"
	httpMiddleware "forum-payment/internal/controller/http/middleware"
)

func Run(
	db model.IDatabase,
	paymentService model.IPaymentService,
	authorizationClient model.IAuthorizationClient,
	serverPort string,
) {
	server := echo.New()
	server.Validator = validator.NewValidator()

	includeHttpMiddleware(server, authorizationClient)
	includeSystemHttpHandler(server, db)
	includePaymentHttpHandler(server, paymentService)

	server.Logger.Fatal(server.Start(":" + serverPort))
}

func includeHttpMiddleware(
	server *echo.Echo,
	authorizationClient model.IAuthorizationClient,
) {
	server.Use(httpMiddleware.LoggerMiddleware())
	server.Use(httpMiddleware.AuthorizationMiddleware(authorizationClient))
}

func includePaymentHttpHandler(
	server *echo.Echo,
	paymentService model.IPaymentService,
) {
	server.POST(model.PREFIX+"/create", paymentHttpHandler.CreatePayment(paymentService))
	server.POST(model.PREFIX+"/paid/:paymentID", paymentHttpHandler.PaidPayment(paymentService))
	server.POST(model.PREFIX+"/cancel/:paymentID", paymentHttpHandler.CancelPayment(paymentService))
	server.GET(model.PREFIX+"/status/:paymentID", paymentHttpHandler.StatusPayment(paymentService))
}

func includeSystemHttpHandler(
	server *echo.Echo,
	db model.IDatabase,
) {
	server.GET(model.PREFIX+"/table/create", CreateTable(db))
	server.GET(model.PREFIX+"/table/drop", DropTable(db))
}

func CreateTable(db model.IDatabase) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := context.Background()

		err := db.CreateTable(ctx, model.CreateTableQuery)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "table created")
	}
}

func DropTable(db model.IDatabase) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := context.Background()

		err := db.DropTable(ctx, model.DropTableQuery)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, "table dropped")
	}
}
