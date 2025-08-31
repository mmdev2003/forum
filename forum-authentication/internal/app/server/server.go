package server

import (
	"context"
	"forum-authentication/internal/model"
	"forum-authentication/pkg/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

import (
	authenticationHttpHandler "forum-authentication/internal/controller/http/handler/authentication"
	httpMiddleware "forum-authentication/internal/controller/http/middleware"
)

func Run(
	db model.IDatabase,
	authenticationService model.IAuthenticationService,
	authorizationClient model.IAuthorizationClient,
	serverPort,
	domain,
	interServerSecretKey string,
) {
	server := echo.New()
	server.Validator = validator.NewValidator()

	includeHttpMiddleware(server, authorizationClient)
	includeSystemHttpHandler(server, db)
	includeAuthenticationHttpHandler(server, authenticationService, domain, interServerSecretKey)

	server.Logger.Fatal(server.Start(":" + serverPort))
}

func includeHttpMiddleware(
	server *echo.Echo,
	authorizationClient model.IAuthorizationClient,
) {
	server.Use(httpMiddleware.LoggerMiddleware())
	server.Use(httpMiddleware.AuthorizationMiddleware(authorizationClient))
}

func includeAuthenticationHttpHandler(
	server *echo.Echo,
	authenticationService model.IAuthenticationService,
	domain,
	interServerSecretKey string,
) {
	server.POST(model.PREFIX+"/register", authenticationHttpHandler.Register(authenticationService, domain))
	server.POST(model.PREFIX+"/login", authenticationHttpHandler.Login(authenticationService, domain))
	server.GET(model.PREFIX+"/2fa/generate", authenticationHttpHandler.GenerateTwoFa(authenticationService))
	server.POST(model.PREFIX+"/2fa/set", authenticationHttpHandler.SetTwoFa(authenticationService))
	server.POST(model.PREFIX+"/2fa/del/:twoFaCode", authenticationHttpHandler.DeleteTwoFaKey(authenticationService))
	server.GET(model.PREFIX+"/2fa/verify/:twoFaCode", authenticationHttpHandler.VerifyTwoFa(authenticationService))
	server.POST(model.PREFIX+"/role/upgrade/admin", authenticationHttpHandler.UpgradeToAdmin(authenticationService, interServerSecretKey))
	server.POST(model.PREFIX+"/role/upgrade/support", authenticationHttpHandler.UpgradeToSupport(authenticationService, interServerSecretKey))
	server.POST(model.PREFIX+"/password/recovery", authenticationHttpHandler.RecoveryPassword(authenticationService))
	server.POST(model.PREFIX+"/password/change", authenticationHttpHandler.ChangePassword(authenticationService))
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
