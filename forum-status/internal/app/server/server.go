package server

import (
	"context"
	"forum-status/internal/model"
	"forum-status/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

import (
	statusHttpHandler "forum-status/internal/controller/http/handler/status"
	httpMiddleware "forum-status/internal/controller/http/middleware"
)

func Run(
	db model.IDatabase,
	statusService model.IStatusService,
	authorizationClient model.IAuthorizationClient,
	interServerSecretKey,
	serverPort string,
) {
	server := echo.New()
	server.Validator = validator.NewValidator()

	includeHttpMiddleware(server, authorizationClient)
	includeSystemHttpHandler(server, db)
	includeStatusHttpHandler(server, statusService, interServerSecretKey)

	server.Logger.Fatal(server.Start(":" + serverPort))
}

func includeHttpMiddleware(
	server *echo.Echo,
	authorizationClient model.IAuthorizationClient,
) {
	server.Use(httpMiddleware.LoggerMiddleware())
	server.Use(httpMiddleware.AuthorizationMiddleware(authorizationClient))
}

func includeStatusHttpHandler(
	server *echo.Echo,
	statusService model.IStatusService,
	interServerSecretKey string,
) {
	server.POST(model.PREFIX+"/payment/create", statusHttpHandler.CreatePaymentForStatus(statusService))
	server.POST(model.PREFIX+"/payment/confirm", statusHttpHandler.ConfirmPaymentForStatus(statusService, interServerSecretKey))
	server.POST(model.PREFIX+"/assign", statusHttpHandler.AssignStatusToAccount(statusService, interServerSecretKey))
	server.POST(model.PREFIX+"/revoke", statusHttpHandler.RevokeStatusFromAccount(statusService, interServerSecretKey))
	server.GET(model.PREFIX+"/info/:accountID", statusHttpHandler.GetStatusByAccountID(statusService))
	server.GET(model.PREFIX+"/all", statusHttpHandler.GetAllStatus(statusService))
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
