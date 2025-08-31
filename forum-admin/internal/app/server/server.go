package server

import (
	"context"
	"forum-admin/internal/model"
	"forum-admin/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

import (
	adminHttpHandler "forum-admin/internal/controller/http/handler/admin"
	httpMiddleware "forum-admin/internal/controller/http/middleware"
)

func Run(
	db model.IDatabase,
	adminService model.IAdminService,
	interServerSecretKey,
	serverPort string,
) {
	server := echo.New()
	server.Validator = validator.NewValidator()

	includeHttpMiddleware(server)
	includeSystemHttpHandler(server, db)
	includeAdminHttpHandler(server, adminService, interServerSecretKey)

	server.Logger.Fatal(server.Start(":" + serverPort))
}

func includeHttpMiddleware(
	server *echo.Echo,
) {
	server.Use(httpMiddleware.LoggerMiddleware())
}

func includeAdminHttpHandler(
	server *echo.Echo,
	adminService model.IAdminService,
	interServerSecretKey string,
) {
	server.POST(model.PREFIX+"/create", adminHttpHandler.CreateAdmin(adminService, interServerSecretKey))
	server.GET(model.PREFIX+"/all", adminHttpHandler.AllAdmin(adminService))
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
