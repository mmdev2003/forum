package server

import (
	"context"
	notificationHandler "forum-notification/internal/api/http/handler/notification"
	"forum-notification/internal/api/http/middleware"
	"forum-notification/internal/api/ws/connection/notification"
	"forum-notification/internal/model"
	"forum-notification/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Run(
	db model.IDatabase,
	notificationService model.INotification,
	notificationFilterService model.INotificationFilter,
	wsConnManager model.IWsConnManager,
	authorizationClient model.IAuthorizationClient,
	serverPort string,
) {
	server := echo.New()
	server.Validator = validator.NewValidator()

	includeWsConnection(server, wsConnManager)
	includeNotificationHandlers(server, notificationService)
	includeHttpMiddleware(server, authorizationClient)
	includeSystemHttpHandler(server, db)
	includeNotificationFilterHandlers(server, notificationFilterService)

	server.Logger.Fatal(server.Start(":" + serverPort))
}

func includeNotificationFilterHandlers(server *echo.Echo, service model.INotificationFilter) {
	server.POST(model.PREFIX+"/filter", notificationHandler.UpdateFiltersByAccountID(service))
	server.GET(model.PREFIX+"/filter", notificationHandler.GetFiltersByAccountID(service))

}

func includeWsConnection(
	server *echo.Echo,
	wsConnManager model.IWsConnManager,
) {
	server.GET(model.PREFIX+"/ws", notification.RegisterNotificationConnection(wsConnManager))
}

func includeNotificationHandlers(
	server *echo.Echo,
	notificationService model.INotification,
) {
	server.POST(model.PREFIX+"/create/:type", notificationHandler.CreateNotification(notificationService))
	server.GET(model.PREFIX+"/", notificationHandler.GetNotificationsByAccountID(notificationService))
}

func includeHttpMiddleware(
	server *echo.Echo,
	authorizationClient model.IAuthorizationClient,
) {
	server.Use(middleware.LoggerMiddleware())
	server.Use(middleware.AuthorizationMiddleware(authorizationClient))
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
