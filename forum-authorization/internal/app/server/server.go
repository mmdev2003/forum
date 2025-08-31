package server

import (
	"context"
	"forum-authorization/internal/model"
	"forum-authorization/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

import (
	authorizationHttpHandler "forum-authorization/internal/controller/http/handler/authorization"
	httpMiddleware "forum-authorization/internal/controller/http/middleware"
)

func Run(
	db model.IDatabase,
	authorizationService model.IAuthorizationService,
	serverPort string,
	domain string,
) {
	server := echo.New()
	server.Validator = validator.NewValidator()

	includeHttpMiddleware(server)
	includeSystemHttpHandler(server, db)
	includeDialogHttpHandler(server, authorizationService, domain)

	server.Logger.Fatal(server.Start(":" + serverPort))
}

func includeHttpMiddleware(
	server *echo.Echo,
) {
	server.Use(httpMiddleware.LoggerMiddleware())
}

func includeDialogHttpHandler(
	server *echo.Echo,
	authorizationService model.IAuthorizationService,
	domain string,
) {
	server.POST(model.PREFIX+"/", authorizationHttpHandler.Authorization(authorizationService))
	server.GET(model.PREFIX+"/check", authorizationHttpHandler.CheckAuthorization(authorizationService))
	server.GET(model.PREFIX+"/refresh", authorizationHttpHandler.RefreshTokens(authorizationService, domain))
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
