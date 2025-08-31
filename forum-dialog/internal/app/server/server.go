package server

import (
	"context"
	"forum-dialog/internal/model"
	"forum-dialog/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

import (
	dialogAmqpHandler "forum-dialog/internal/controller/amqp/handler/dialog"
	dialogHttpHandler "forum-dialog/internal/controller/http/handler/dialog"
	httpMiddleware "forum-dialog/internal/controller/http/middleware"
	dialogWsHandler "forum-dialog/internal/controller/ws/connection/dialog"
)

func Run(
	db model.IDatabase,
	dialogService model.IDialogService,
	authorizationClient model.IAuthorizationClient,
	wsConnManager model.IWsConnManager,
	messageBroker model.IMessageBroker,
	serverPort string,
	podID string,
) {
	server := echo.New()
	server.Validator = validator.NewValidator()

	includeHttpMiddleware(server, authorizationClient)
	includeSystemHttpHandler(server, db)
	includeDialogWsConnection(server, dialogService, wsConnManager)
	includeDialogHttpHandler(server, dialogService)
	includeDialogAmqpHandler(dialogService, messageBroker, podID)

	server.Logger.Fatal(server.Start(":" + serverPort))
}

func includeHttpMiddleware(
	server *echo.Echo,
	authorizationClient model.IAuthorizationClient,
) {
	server.Use(httpMiddleware.LoggerMiddleware())
	server.Use(httpMiddleware.AuthorizationMiddleware(authorizationClient))
}

func includeDialogAmqpHandler(
	dialogService model.IDialogService,
	messageBroker model.IMessageBroker,
	podID string,
) {
	go messageBroker.Subscribe(podID, dialogAmqpHandler.SendMessageToDialog(dialogService))
}

func includeDialogHttpHandler(
	server *echo.Echo,
	dialogService model.IDialogService,
) {
	server.POST(model.PREFIX+"/create", dialogHttpHandler.CreateDialog(dialogService))
	server.POST(model.PREFIX+"/mark/star/:dialogID", dialogHttpHandler.MarkDialogAsStarred(dialogService))
	server.POST(model.PREFIX+"/mark/read/:dialogID", dialogHttpHandler.MarkMessagesAsRead(dialogService))
	server.POST(model.PREFIX+"/del/:dialogID", dialogHttpHandler.DeleteDialog(dialogService))
	server.GET(model.PREFIX+"/all", dialogHttpHandler.DialogsByAccountID(dialogService))
	server.GET(model.PREFIX+"/messages/:dialogID", dialogHttpHandler.MessagesByDialogID(dialogService))
	server.POST(model.PREFIX+"/file/upload", dialogHttpHandler.UploadFile(dialogService))
	server.GET(model.PREFIX+"/file/download/:fileURL", dialogHttpHandler.DownloadFile(dialogService))
}

func includeDialogWsConnection(
	server *echo.Echo,
	dialogService model.IDialogService,
	wsConnManager model.IWsConnManager,
) {
	server.GET(model.PREFIX+"/ws", dialogWsHandler.RegisterDialogConnection(dialogService, wsConnManager))
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
