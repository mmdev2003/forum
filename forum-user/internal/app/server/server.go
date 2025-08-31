package server

import (
	"context"

	"forum-user/internal/model"
	"forum-user/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

import (
	userHttpHandler "forum-user/internal/controller/http/handler/user"
	httpMiddleware "forum-user/internal/controller/http/middleware"
)

func Run(
	db model.IDatabase,
	fullTextSearchEngine model.IFullTextSearchEngine,
	userService model.IUserService,
	authorizationClient model.IAuthorizationClient,
	serverPort string,
) {
	server := echo.New()
	server.Validator = validator.New()

	includeSystemHttpHandler(server, db, fullTextSearchEngine)
	includeHttpMiddleware(server, authorizationClient)
	includeUserHttpHandler(server, userService)

	server.Logger.Fatal(server.Start(":" + serverPort))
}

func includeUserHttpHandler(
	server *echo.Echo,
	userService model.IUserService,
) {
	server.POST(model.PREFIX+"/create", userHttpHandler.CreateUser(userService))
	server.POST(model.PREFIX+"/ban/create", userHttpHandler.BanUser(userService))
	server.POST(model.PREFIX+"/ban/del/:toAccountID", userHttpHandler.UnbanUser(userService))
	server.POST(model.PREFIX+"/warning/create", userHttpHandler.NewWarningFromAdmin(userService))
	server.POST(model.PREFIX+"/avatar/upload", userHttpHandler.UploadAvatar(userService))
	server.POST(model.PREFIX+"/header/upload", userHttpHandler.UploadHeader(userService))
	server.GET(model.PREFIX+"/avatar/download/:fileID", userHttpHandler.DownloadAvatar(userService))
	server.GET(model.PREFIX+"/header/download/:fileID", userHttpHandler.DownloadHeader(userService))
	server.GET(model.PREFIX+"/info/id/:accountID", userHttpHandler.GetUserByAccountID(userService))
	server.GET(model.PREFIX+"/search/login/:login", userHttpHandler.UsersByLoginSearch(userService))
	server.GET(model.PREFIX+"/info/login/:login", userHttpHandler.GetUserByLogin(userService))
	server.GET(model.PREFIX+"/info/bans", userHttpHandler.BanByAccountID(userService))
	server.GET(model.PREFIX+"/warning/all/:toAccountID", userHttpHandler.AllWarningFromAdmin(userService))
}

func includeHttpMiddleware(
	server *echo.Echo,
	authorizationClient model.IAuthorizationClient,
) {
	server.Use(httpMiddleware.LoggerMiddleware())
	server.Use(httpMiddleware.AuthorizationMiddleware(authorizationClient))
}

func includeSystemHttpHandler(
	server *echo.Echo,
	db model.IDatabase,
	fullTextSearchEngine model.IFullTextSearchEngine,
) {
	server.GET(model.PREFIX+"/table/create", createTable(db))
	server.GET(model.PREFIX+"/table/drop", dropTable(db))
	server.GET(model.PREFIX+"/fullTextSearchEngine/create", createFullTextSearchEngineIndex(fullTextSearchEngine))
}

func createFullTextSearchEngineIndex(fullTextSearchEngine model.IFullTextSearchEngine) echo.HandlerFunc {
	return func(request echo.Context) error {
		indexes := []string{
			model.UserFullTextSearchIndex,
		}
		err := fullTextSearchEngine.CreateIndexes(indexes)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "indexes created")
	}
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
