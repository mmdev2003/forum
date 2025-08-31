package server

import (
	"context"
	"forum-frame/internal/model"
	"forum-frame/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

import (
	frameHttpHandler "forum-frame/internal/controller/http/handler/frame"
	httpMiddleware "forum-frame/internal/controller/http/middleware"
)

func Run(
	db model.IDatabase,
	frameService model.IFrameService,
	authorizationClient model.IAuthorizationClient,
	interServerSecretKey,
	serverPort string,
) {
	server := echo.New()
	server.Validator = validator.NewValidator()

	includeHttpMiddleware(server, authorizationClient)
	includeSystemHttpHandler(server, db)
	includeFrameHttpHandler(server, frameService, interServerSecretKey)

	server.Logger.Fatal(server.Start(":" + serverPort))
}

func includeHttpMiddleware(
	server *echo.Echo,
	authorizationClient model.IAuthorizationClient,
) {
	server.Use(httpMiddleware.LoggerMiddleware())
	server.Use(httpMiddleware.AuthorizationMiddleware(authorizationClient))
}

func includeFrameHttpHandler(
	server *echo.Echo,
	frameService model.IFrameService,
	interServerSecretKey string,
) {
	server.POST(model.PREFIX+"/payment/create", frameHttpHandler.CreatePaymentForFrame(frameService))
	server.POST(model.PREFIX+"/payment/confirm", frameHttpHandler.ConfirmPaymentForFrame(frameService, interServerSecretKey))
	server.POST(model.PREFIX+"/new", frameHttpHandler.AddNewFrame(frameService))
	server.POST(model.PREFIX+"/change/:dbFrameID", frameHttpHandler.ChangeCurrentFrame(frameService))
	server.GET(model.PREFIX+"/:accountID", frameHttpHandler.FramesByAccountID(frameService))
	server.GET(model.PREFIX+"/all", frameHttpHandler.AllFrame(frameService))
	server.GET(model.PREFIX+"/download/:frameID", frameHttpHandler.DownloadFrame(frameService))
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
