package server

import (
	"context"
	dialogHttpHandler "forum-support/internal/controller/http/handler/dialog"
	supportHttpHandler "forum-support/internal/controller/http/handler/support-request"
	httpMiddleware "forum-support/internal/controller/http/middleware"
	dialogWsHandler "forum-support/internal/controller/ws/connection/dialog"
	"forum-support/internal/model"
	"forum-support/pkg/validator"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Run(
	db model.IDatabase,
	supportRequestService model.ISupportRequestService,
	dialogService model.IDialogService,
	authorizationClient model.IAuthorizationClient,
	wsConnManager model.IWsConnManager,
	serverPort string,
) {
	server := echo.New()
	server.Validator = validator.NewValidator()

	includeSystemHttpHandler(server, db)
	includeSupportRequestHttpHandlers(server, supportRequestService)
	includeDialogHttpHandlers(server, dialogService)
	includeHttpMiddleware(server, authorizationClient)
	includeDialogWsConnection(server, dialogService, wsConnManager)

	server.Logger.Fatal(server.Start(":" + serverPort))
}

func includeDialogHttpHandlers(server *echo.Echo, dialogService model.IDialogService) {
	server.POST(model.PREFIX+"/dialog/create", dialogHttpHandler.CreateDialog(dialogService))
	server.POST(model.PREFIX+"/dialog/:dialogID/mark/read", dialogHttpHandler.MarkMessagesAsRead(dialogService))
	server.GET(model.PREFIX+"/dialog/:dialogID", dialogHttpHandler.MessagesByDialogID(dialogService))
	server.GET(model.PREFIX+"/dialog/all", dialogHttpHandler.AllDialogs(dialogService))
}

func includeDialogWsConnection(
	server *echo.Echo,
	dialogService model.IDialogService,
	wsConnManager model.IWsConnManager,
) {
	server.GET(model.PREFIX+"/ws", dialogWsHandler.RegisterDialogConnection(dialogService, wsConnManager))
}

func includeHttpMiddleware(
	server *echo.Echo,
	authorizationClient model.IAuthorizationClient,
) {
	server.Use(httpMiddleware.LoggerMiddleware())
	server.Use(httpMiddleware.AuthorizationMiddleware(authorizationClient))
}

func includeSupportRequestHttpHandlers(
	server *echo.Echo,
	supportRequestService model.ISupportRequestService,
) {
	server.GET(model.PREFIX+"/request", supportHttpHandler.GetRequests(supportRequestService))
	server.GET(model.PREFIX+"/request/:requestID", supportHttpHandler.GetRequest(supportRequestService))
	server.POST(model.PREFIX+"/request/create", supportHttpHandler.CreateRequest(supportRequestService))
	server.POST(model.PREFIX+"/request/open/:requestID", supportHttpHandler.OpenRequest(supportRequestService))
	server.POST(model.PREFIX+"/request/close/:requestID", supportHttpHandler.CloseRequest(supportRequestService))
}

func includeSystemHttpHandler(
	server *echo.Echo,
	db model.IDatabase,
) {
	server.GET(model.PREFIX+"/table/create", createTable(db))
	server.GET(model.PREFIX+"/table/drop", dropTable(db))
}

func createTable(db model.IDatabase) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := context.Background()

		err := db.CreateTable(ctx, model.CreateTableQuery)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "table created")
	}
}

func dropTable(db model.IDatabase) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := context.Background()

		err := db.DropTable(ctx, model.DropTableQuery)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, "table dropped")
	}
}
